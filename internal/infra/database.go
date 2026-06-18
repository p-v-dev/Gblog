package infra

import (
	"Gblog/internal/blogPost"
	"Gblog/internal/comment"
	"Gblog/internal/tag"
	"Gblog/internal/user"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadEnv() {
	// ponytail: loads .env from CWD, go run is always from project root
	_ = godotenv.Load()
}

func ConnectDB() (*gorm.DB, error) {
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}
	dsn := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") + " port=" + os.Getenv("DB_PORT") + " sslmode=" + sslmode + " TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&blogPost.BlogPost{}, &user.User{}, &tag.Tag{}, &comment.Comment{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
