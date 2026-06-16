package http

import (
	"Gblog/internal/blogPost"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BlogPostHandler struct {
	createUseCase  blogPost.CreatePostUseCase
	publishUseCase blogPost.PublishPostUseCase
	updateUseCase  blogPost.UpdatePostUseCase
	deleteUseCase  blogPost.DeletePostUseCase
}

func NewBlogPostHandler(
	create blogPost.CreatePostUseCase,
	publish blogPost.PublishPostUseCase,
	update blogPost.UpdatePostUseCase,
	delete blogPost.DeletePostUseCase,
) *BlogPostHandler {
	return &BlogPostHandler{
		createUseCase:  create,
		publishUseCase: publish,
		updateUseCase:  update,
		deleteUseCase:  delete,
	}
}

// Create godoc
// @Summary      Criar um Blog Post
// @Description  Cria um novo post no blog salvando-o inicialmente como rascunho (draft).
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        post  body      blogPost.CreatePostInputDTO  true  "Dados para a criação do Post"
// @Success      201   {object}  map[string]string            "Post criado com sucesso"
// @Failure      400   {object}  map[string]string            "JSON inválido"
// @Failure      422   {object}  map[string]string            "Erro de regra de negócio"
// @Router       /posts [post]
func (h *BlogPostHandler) Create(c *gin.Context) {
	var input blogPost.CreatePostInputDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
		return
	}

	if err := h.createUseCase.Execute(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Post criado com sucesso como rascunho!"})
}

// Publish godoc
// @Summary      Publicar um Blog Post
// @Description  Altera o status de um post existente para 'published' se ele cumprir os requisitos de negócio.
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        id    path      string               true  "ID do Post (UUID)"
// @Success      200   {object}  map[string]string  "Post publicado com sucesso"
// @Failure      400   {object}  map[string]string  "ID inválido ou erro na publicação"
// @Router       /posts/{id}/publish [patch]
func (h *BlogPostHandler) Publish(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido, deve ser um UUID válido"})
		return
	}

	if err := h.publishUseCase.Execute(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post publicado com sucesso!"})
}

// Update godoc
// @Summary      Editar um Blog Post
// @Description  Atualiza o título, slug e conteúdo de um post ativo.
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        id    path      string                       true  "ID do Post (UUID)"
// @Param        post  body      blogPost.UpdatePostInputDTO  true  "Novos dados do Post"
// @Success      200   {object}  map[string]string            "Post atualizado com sucesso"
// @Failure      400   {object}  map[string]string            "ID ou JSON inválido"
// @Router       /posts/{id} [put]
func (h *BlogPostHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido, deve ser um UUID válido"})
		return
	}

	var input blogPost.UpdatePostInputDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	if err := h.updateUseCase.Execute(c.Request.Context(), id, input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post atualizado com sucesso!"})
}

// Delete godoc
// @Summary      Deletar um Blog Post
// @Description  Inativa um post do sistema realizando uma exclusão lógica (soft delete).
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        id    path      string               true  "ID do Post (UUID)"
// @Success      200   {object}  map[string]string  "Post removido com sucesso"
// @Failure      400   {object}  map[string]string  "ID inválido ou post não encontrado"
// @Router       /posts/{id} [delete]
func (h *BlogPostHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido, deve ser um UUID válido"})
		return
	}

	if err := h.deleteUseCase.Execute(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post removido com sucesso (soft delete)!"})
}
