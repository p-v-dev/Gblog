package infra

import (
	"Gblog/internal/blogPost"
	"os"

	"log"

	"path/filepath"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadEnv() {
	// 1. Pega o diretório atual de execução
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// 2. Sobe as pastas até encontrar o arquivo go.mod (que define a raiz do projeto)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break // Encontrou a raiz!
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Chegou no topo do sistema de arquivos e não achou a raiz
			break
		}
		dir = parent
	}

	// 3. Carrega o .env a partir da raiz encontrada
	envPath := filepath.Join(dir, ".env")

	// Use Load se o .env for obrigatório, ou Overload se quiser sobrescrever variáveis do sistema
	_ = godotenv.Load(envPath)
}

func ConnectDB() (*gorm.DB, error) {
	dsn := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") + " port=" + os.Getenv("DB_PORT") + " sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&blogPost.BlogPost{}) // Certifique-se de importar o pacote onde está o BlogPost se necessário
	if err != nil {
		return nil, err
	}
	return db, nil
}
