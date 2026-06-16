package main

import (
	"Gblog/internal/blogPost" // <-- Puxa tanto a struct quanto o Repository
	blogHttp "Gblog/internal/blogPost/http"
	"Gblog/internal/infra"
	"Gblog/internal/user"
	userHttp "Gblog/internal/user/http"
	"log"
	"os"

	// IMPORTANTE: Importa os documentos que serão gerados pelo comando 'swag init'
	_ "Gblog/cmd/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Gblog API
// @version         1.0
// @description     API de gerenciamento de artigos de um blog usando Clean Architecture.
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	infra.LoadEnv()
	db, err := infra.ConnectDB()
	if err != nil {
		log.Fatalf("Erro no banco: %v", err)
	}
	db.AutoMigrate(&blogPost.BlogPost{})

	repo := blogPost.NewBlogPostRepository(db)
	createUseCase := blogPost.NewCreatePostUseCase(repo)
	publishUseCase := blogPost.NewPublishPostUseCase(repo)
	updateUseCase := blogPost.NewUpdatePostUseCase(repo)
	deleteUseCase := blogPost.NewDeletePostUseCase(repo)

	r := gin.Default()

	// ROTA DO SWAGGER: Configura a rota onde a interface gráfica do Swagger vai rodar
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	postHandler := blogHttp.NewBlogPostHandler(
		createUseCase,
		publishUseCase,
		updateUseCase,
		deleteUseCase,
	)

	userRepo := user.NewUserRepository(db)
	createUserUseCase := user.NewCreateUserUseCase(userRepo)
	userHandler := userHttp.NewUserHandler(createUserUseCase)

	api := r.Group("/api/v1")
	{
		api.POST("/posts", postHandler.Create)               // Criar rascunho
		api.PUT("/posts/:id", postHandler.Update)            // Editar conteúdo
		api.PATCH("/posts/:id/publish", postHandler.Publish) // Ação específica de publicar
		api.DELETE("/posts/:id", postHandler.Delete)         // Soft Delete
		api.POST("/users", userHandler.Create)               // Criar usuário
	}

	r.Run(":" + os.Getenv("API_PORT"))

}
