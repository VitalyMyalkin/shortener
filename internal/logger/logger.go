package logger

import (
    "time"
	"net/http"


	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Log будет доступен всему коду как синглтон.
// Никакой код навыка, кроме функции InitLogger, не должен модифицировать эту переменную.
// По умолчанию установлен no-op-логер, который не выводит никаких сообщений.
var Log *zap.Logger = zap.NewNop()

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
func Initialize() error {
    // создаём новую конфигурацию логера
    cfg := zap.NewProductionConfig()
    // создаём логер на основе конфигурации
    zl, err := cfg.Build()
    if err != nil {
        return err
    }
    // устанавливаем синглтон
    Log = zl
    return nil
}

type (
    // берём структуру для хранения сведений об ответе
    responseData struct {
        size int
    }

    // добавляем реализацию http.ResponseWriter
    loggingResponseWriter struct {
        http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
        responseData *responseData
    }
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
    // записываем ответ, используя оригинальный http.ResponseWriter
    size, err := r.ResponseWriter.Write(b) 
    r.responseData.size += size // захватываем размер
    return size, err
}

// WithLogging добавляет дополнительный код для регистрации сведений о запросе
// и возвращает новый http.Handler.
func WithLogging() gin.HandlerFunc {
    logFn := func(c *gin.Context) {
        start := time.Now()

        responseData := &responseData {
            size: 0,
        }
        
        c.Next()

        duration := time.Since(start)

        Log.Info("got incoming HTTP request "+c.Request.URL.Path,
            zap.String("uri", c.Request.RequestURI),
            zap.String("method", c.Request.Method),
            zap.Int("status", c.Writer.Status()), // получаем код статуса ответа
            zap.Duration("duration", duration),
            zap.Int("size", responseData.size), // получаем перехваченный размер ответа
        )
    }
    return logFn
} 

