package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

func (g *gzipWriter) Close() error {
	return g.writer.Close()
}

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := c.Request.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			gz, err := gzip.NewWriterLevel(c.Writer, gzip.DefaultCompression)
			if err != nil {
				io.WriteString(c.Writer, err.Error())
				return
			}
			c.Writer = &gzipWriter{c.Writer, gz}
			c.Header("Content-Encoding", "gzip")
			c.Header("Vary", "Accept-Encoding")

			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer gz.Close()
			c.Next()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := c.Request.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := newCompressReader(c.Request.Body)
			if err != nil {
				c.Writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			c.Request.Body = cr
			defer cr.Close()
			c.Next()
		}
	}
}
