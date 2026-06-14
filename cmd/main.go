package main

import (
	"context"
	"log"

	"Gblog/internal/blogPost" // <-- Puxa tanto a struct quanto o Repository
	"Gblog/internal/infra"
)

func main() {
	infra.LoadEnv()

	db, err := infra.ConnectDB()
	if err != nil {
		log.Fatalf("Erro no banco: %v", err)
	}

	// O GORM precisa conhecer a struct para a migração automática
	// Como você mudou a struct de lugar, ajuste a linha abaixo no seu ConnectDB ou chame aqui:
	db.AutoMigrate(&blogPost.BlogPost{})

	// Instancia o repositório usando o pacote blogPost
	postRepo := blogPost.NewBlogPostRepository(db)
	postUseCase := blogPost.NewCreatePostUseCase(postRepo)
	ctx := context.Background()

	// E

	err = postUseCase.Execute(ctx, blogPost.CreatePostInputDTO{
		Title:   "Arquitetura Limpa em Go",
		Slug:    "arquitetura-limpa-em-go",
		Content: "Isolando suas entidades com DTOs e Interfaces!",
	})
	if err != nil {
		log.Printf("Erro ao criar post pelo UseCase: %v", err)
	} else {
		log.Println("Post criado com sucesso através do Use Case!")
	}

	log.Println("Post criado com sucesso!")
}
