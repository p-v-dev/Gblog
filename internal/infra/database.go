package infra

import (
	"Gblog/internal/blogPost"
	"Gblog/internal/comment"
	"Gblog/internal/tag"
	"Gblog/internal/user"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadEnv() {
	// ponytail: loads .env from CWD, go run is always from project root
	_ = godotenv.Load()
}

func ConnectDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		tz := os.Getenv("DB_TZ")
		if tz == "" {
			tz = "UTC"
		}
		sslmode := os.Getenv("DB_SSLMODE")
		if sslmode == "" {
			sslmode = "disable"
		}
		dsn = "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") + " port=" + os.Getenv("DB_PORT") + " sslmode=" + sslmode + " TimeZone=" + tz
	}

	// ponytail: simple retry, add exponential backoff if DB takes >10s to start
	var db *gorm.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		slog.Warn("tentativa de conexão falhou, tentando novamente...", "tentativa", i+1, "err", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar após 5 tentativas: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	err = db.AutoMigrate(&blogPost.BlogPost{}, &user.User{}, &tag.Tag{}, &comment.Comment{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
