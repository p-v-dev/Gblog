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
	postUseCase := blogPost.NewBlogPostUseCase(postRepo)
	ctx := context.Background()

	// Exemplo de uso
	novoPostDTO := blogPost.BlogPostInputDTO{
		Title:   "Arquitetura Limpa em Go",
		Slug:    "arquitetura-limpa-em-go",
		Content: "Isolando suas entidades com DTOs e Interfaces!",
		Status:  "draft", // Lembra que validamos no Use Case? Tem que ser draft, published ou archived
	}

	err = postUseCase.Create(ctx, novoPostDTO)
	if err != nil {
		log.Printf("Erro ao criar post pelo UseCase: %v", err)
	} else {
		log.Println("Post criado com sucesso através do Use Case!")
	}

	log.Println("Post criado com sucesso!")
}
