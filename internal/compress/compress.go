package compress

import (
    "strings"
	"compress/gzip"
    "io"
    "net/http"

	"github.com/gin-gonic/gin"
)

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type compressWriter struct {
    w  http.ResponseWriter
    zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
    return &compressWriter{
        w:  w,
        zw: gzip.NewWriter(w),
    }
}

func (c *compressWriter) Header() http.Header {
    return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
    return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
    if statusCode < 300 {
        c.w.Header().Set("Content-Encoding", "gzip")
    }
    c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
    return c.zw.Close()
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
            cw := newCompressWriter(c.Writer)
            // не забываем отправить клиенту все сжатые данные после завершения middleware
            defer cw.Close()
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
		}
        c.Next()
	}
}