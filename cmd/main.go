package main

import (
	"Gblog/internal/blogPost" // <-- Puxa tanto a struct quanto o Repository
	blogHttp "Gblog/internal/blogPost/http"
	"Gblog/internal/comment"
	commentHttp "Gblog/internal/comment/http"
	"Gblog/internal/infra"
	"Gblog/internal/tag"
	tagHttp "Gblog/internal/tag/http"
	"Gblog/internal/user"
	userHttp "Gblog/internal/user/http"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// IMPORTANTE: Importa os documentos que serão gerados pelo comando 'swag init'
	_ "Gblog/cmd/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Gblog API
// @version         1.0
// @description     API de gerenciamento de artigos de um blog usando Clean Architecture.
// @host            localhost:8080
// @BasePath        /api/v1
// ponytail: adapter so user.Repository satisfies blogPost.UserExistenceChecker
type userExistenceChecker struct {
	repo user.Repository
}

func (c *userExistenceChecker) UserExists(ctx context.Context, userID string) bool {
	_, err := c.repo.FindByID(ctx, userID)
	return err == nil
}

// ponytail: adapter so blogPost.BlogPostRepository satisfies comment.PostExistenceChecker
type postExistenceChecker struct {
	repo blogPost.BlogPostRepository
}

func (c *postExistenceChecker) PostExists(ctx context.Context, id string) bool {
	_, err := c.repo.FindByID(ctx, id)
	return err == nil
}

type authHandler struct {
	loginUseCase user.LoginUseCase
}

func (h *authHandler) Login(c *gin.Context) {
	var body user.LoginInput
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email e senha são obrigatórios"})
		return
	}

	output, err := h.loginUseCase.Execute(c.Request.Context(), body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := infra.GenerateToken(output.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao gerar token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func main() {
	infra.LoadEnv()
	db, err := infra.ConnectDB()
	if err != nil {
		log.Fatalf("Erro no banco: %v", err)
	}
	repo := blogPost.NewBlogPostRepository(db)
	userRepo := user.NewUserRepository(db)
	tagRepo := tag.NewTagRepository(db)
	createUseCase := blogPost.NewCreatePostUseCase(repo, &userExistenceChecker{repo: userRepo}, tagRepo)
	publishUseCase := blogPost.NewPublishPostUseCase(repo)
	updateUseCase := blogPost.NewUpdatePostUseCase(repo, tagRepo)
	deleteUseCase := blogPost.NewDeletePostUseCase(repo)
	getPostsUseCase := blogPost.NewGetPostsUseCase(repo)
	getPostBySlugUseCase := blogPost.NewGetPostBySlugUseCase(repo)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// ROTA DO SWAGGER: Configura a rota onde a interface gráfica do Swagger vai rodar
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	postHandler := blogHttp.NewBlogPostHandler(
		createUseCase,
		publishUseCase,
		updateUseCase,
		deleteUseCase,
		getPostsUseCase,
		getPostBySlugUseCase,
	)

	createUserUseCase := user.NewCreateUserUseCase(userRepo)
	getUserUseCase := user.NewGetUserUseCase(userRepo)
	updateUserUseCase := user.NewUpdateUserUseCase(userRepo)
	deleteUserUseCase := user.NewDeleteUserUseCase(userRepo)
	loginUseCase := user.NewLoginUseCase(userRepo)
	userHandler := userHttp.NewUserHandler(createUserUseCase, getUserUseCase, updateUserUseCase, deleteUserUseCase)

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET não configurado")
	}

	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		log.Fatal("API_PORT não configurada")
	}

	auth := &authHandler{loginUseCase: loginUseCase}

	createTagUseCase := tag.NewCreateTagUseCase(tagRepo)
	listTagsUseCase := tag.NewListTagsUseCase(tagRepo)
	tagHandler := tagHttp.NewTagHandler(createTagUseCase, listTagsUseCase)

	commentRepo := comment.NewCommentRepository(db)
	createCommentUseCase := comment.NewCreateCommentUseCase(commentRepo, &postExistenceChecker{repo: repo})
	listCommentsUseCase := comment.NewListCommentsUseCase(commentRepo)
	deleteCommentUseCase := comment.NewDeleteCommentUseCase(commentRepo)
	commentHandler := commentHttp.NewCommentHandler(createCommentUseCase, listCommentsUseCase, deleteCommentUseCase)

	api := r.Group("/api/v1")
	{
		api.POST("/auth/token", auth.Login) // Gerar token

		api.GET("/posts", postHandler.GetPosts)                 // Listar posts
		api.GET("/posts/slug/:slug", postHandler.GetPostBySlug) // Ver post por slug
		api.GET("/posts/:id/comments", commentHandler.List)     // Listar comentários

		protected := api.Group("")
		protected.Use(infra.AuthMiddleware())
		{
			protected.POST("/posts", postHandler.Create)                 // Criar rascunho
			protected.POST("/posts/:id/comments", commentHandler.Create) // Comentar
			protected.PUT("/posts/:id", postHandler.Update)              // Editar conteúdo
			protected.PATCH("/posts/:id/publish", postHandler.Publish)   // Publicar
			protected.DELETE("/posts/:id", postHandler.Delete)           // Soft Delete
			protected.POST("/tags", tagHandler.Create)                   // Criar tag
			protected.DELETE("/comments/:id", commentHandler.Delete)     // Deletar comentário
			protected.PUT("/users/:id", userHandler.Update)              // Atualizar usuário
			protected.DELETE("/users/:id", userHandler.Delete)           // Desativar usuário
		}

		api.GET("/tags", tagHandler.List)      // Listar tags
		api.GET("/users/:id", userHandler.Get) // Ver usuário
		api.POST("/users", userHandler.Create) // Criar usuário
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Erro ao obter conexão SQL: %v", err)
	}

	srv := &http.Server{
		Addr:         ":" + apiPort,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Println("Desligando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Erro ao desligar servidor: %v", err)
	}

	sqlDB.Close()

}
