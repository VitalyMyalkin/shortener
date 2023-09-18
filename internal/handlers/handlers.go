package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"errors"
	"bufio"
	"context"
	"database/sql"
	"time"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/google/uuid"

	"github.com/VitalyMyalkin/shortener/internal/config"
	"github.com/VitalyMyalkin/shortener/internal/logger"
	"github.com/VitalyMyalkin/shortener/internal/storage"
)

type App struct {
	Cfg     config.Config
	Storage *storage.Storage
	PostgresDB *sql.DB		
}

type Request struct {
	URLstring string `json:"url"`
}

type URLwithID struct {
	ID string `json:"correlation_id"`
	URL string `json:"original_url"`
	Shortened string `json:"short_url"`
}

func NewApp() *App {
	cfg := config.GetConfig()
	storage := storage.NewStorage()
	db, err := sql.Open("pgx", cfg.PostgresDBAddr)
    if err != nil {
        fmt.Println(err)
    }
	return &App{
		Cfg:     cfg,
		Storage: storage,
		PostgresDB: db,
	}
}

func (newApp *App) PingPostgresDB(c *gin.Context) {
	if err := newApp.PostgresDB.PingContext(context.Background()); err != nil {
        c.Status(http.StatusInternalServerError)
    } else {
		c.Status(http.StatusOK)
	}
}

func (newApp *App) AddOrigin(url *url.URL) (string, error) {
	var err error
	short := uuid.NewString()

	if newApp.Cfg.PostgresDBAddr != "" {
		_, err := newApp.PostgresDB.Exec("CREATE TABLE IF NOT EXISTS urls (id SERIAL PRIMARY KEY, origin TEXT UNIQUE, shortened TEXT)")
		if err != nil {
			logger.Log.Fatal(err.Error())
		}
		_, err = newApp.PostgresDB.Exec("INSERT INTO urls (origin, shortened) VALUES ($1, $2)", url.String(), short)
		if err != nil {
			var pgErr *pgconn.PgError
        	if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
				var a string
				newApp.PostgresDB.QueryRow("SELECT shortened FROM urls WHERE origin = $1", url.String()).Scan(&a)
				return a, err
			}
			logger.Log.Fatal("возникла другая ошибка при записи ссылки в бд")
		}
	} else if newApp.Cfg.FilePath != "" {
		fileName := newApp.Cfg.FilePath
		Producer, err := storage.NewFileWriter(fileName)
		if err != nil {
			logger.Log.Fatal("не создан или не открылся файл записи" + fileName)
		}
		defer Producer.Close()
		if err := Producer.WriteShortenedURL(short, url); err != nil {
			logger.Log.Fatal("запись не внесена в файл" + fileName)
		}
	} else {
		newApp.Storage.AddOrigin(short, url)
	}
	return short, err
}


func (newApp *App) GetShortened(c *gin.Context) {
	var pgErr *pgconn.PgError
	var result string
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
	url, err := url.ParseRequestURI(string(body))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": string(body) + "не является валидным URL",
		})
	}

	result, err = newApp.AddOrigin(url)

	if err == nil {
		c.Header("Content-Type", "text/plain")
		c.String(http.StatusCreated, newApp.Cfg.ShortenAddr+"/"+result)
	} else if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
		c.Header("Content-Type", "text/plain")
		c.String(http.StatusConflict, newApp.Cfg.ShortenAddr+"/"+result)
	}
}

func (newApp *App) SendBatch (c *gin.Context) {
	// десериализуем запрос 
	logger.Log.Debug("decoding request")
	var urls []URLwithID
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	json.Unmarshal([]byte(body), &urls)

	if newApp.Cfg.PostgresDBAddr != "" {
		// начинаю транзакцию
		tx, err := newApp.PostgresDB.BeginTx(context.Background(), nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
		}
		
		for i := 0; i < len(urls); i++ {
			if urls[i].ID !="" && urls[i].URL !="" {
				short := uuid.NewString()
				// все изменения записываются в транзакцию
				_, err := tx.ExecContext(context.Background(), 
				"INSERT INTO urls (origin, shortened) VALUES ($1, $2)", urls[i].URL, short)
				if err != nil {
					// если ошибка, то откатываем изменения
					tx.Rollback()
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
				}
				urls[i].Shortened = newApp.Cfg.ShortenAddr + "/" + short
			}
		}

		// коммитим транзакцию
		tx.Commit()
	} else {
		for i := 0; i < len(urls); i++  {
			if urls[i].ID !="" && urls[i].URL !="" {
				url, err := url.ParseRequestURI(urls[i].URL)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": string(body) + "не является валидным URL",
					})
				}
				result, erro := newApp.AddOrigin(url)
				if erro != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
				}
				urls[i].Shortened = newApp.Cfg.ShortenAddr + "/" + result
			}
		}
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, &urls)
}

func (newApp *App) GetShortenedAPI (c *gin.Context) {
	var pgErr *pgconn.PgError
	var result string
	// десериализуем запрос в структуру модели
	logger.Log.Debug("decoding request")
	var req Request
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		c.String(http.StatusInternalServerError, "")
	}

	url, err := url.ParseRequestURI(req.URLstring)
	if err != nil {
		logger.Log.Debug(req.URLstring+"не является валидным URL", zap.Error(err))
		c.String(http.StatusBadRequest, "")
	}

	result, err = newApp.AddOrigin(url)

	if err == nil {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusCreated, gin.H{
			"result": newApp.Cfg.ShortenAddr + "/" + result,
		})
	} else if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusConflict, gin.H{
			"result": newApp.Cfg.ShortenAddr + "/" + result,
		})
	}
}

func (newApp *App) GetOrigin(c *gin.Context) {
	var original string

	if newApp.Cfg.PostgresDBAddr != "" {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
    	defer cancel()
    	// делаем обращение к db в рамках полученного контекста
    	row := newApp.PostgresDB.QueryRowContext(ctx, "SELECT origin FROM urls WHERE shortened = $1", c.Param("id"))
    	// готовим переменную для чтения результата
    	err := row.Scan(&original)  // разбираем результат
    	if err != nil {
        	logger.Log.Fatal("невозможно прочитать данные записи из базы данных")
    	}
	} else if newApp.Cfg.FilePath != "" {
		fileName := newApp.Cfg.FilePath

		file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			logger.Log.Fatal("не создан или не открылся файл записи")
		}

		if err != nil {
			logger.Log.Fatal("невозможно прочитать данные файла записи")
		}

		var shortenedURL storage.ShortenedURL

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			err = json.Unmarshal(scanner.Bytes(), &shortenedURL)
			if err != nil {
				logger.Log.Fatal("не создана структура")
			} else {
				if shortenedURL.ShortURL == c.Param("id") {
					original = shortenedURL.OriginalURL
				}
			}
		}
	} else {
		original = newApp.Storage.Storage[c.Param("id")]
	}

	c.Header("Location", original)
	c.Status(http.StatusTemporaryRedirect)
}
