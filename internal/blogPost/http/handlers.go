package http

import (
	"Gblog/internal/blogPost"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BlogPostHandler struct {
	createUseCase       blogPost.CreatePostUseCase
	publishUseCase      blogPost.PublishPostUseCase
	updateUseCase       blogPost.UpdatePostUseCase
	deleteUseCase       blogPost.DeletePostUseCase
	getPostsUseCase     blogPost.GetPostsUseCase
	getPostBySlugUseCase blogPost.GetPostBySlugUseCase
}

func NewBlogPostHandler(
	create blogPost.CreatePostUseCase,
	publish blogPost.PublishPostUseCase,
	update blogPost.UpdatePostUseCase,
	delete blogPost.DeletePostUseCase,
	getPosts blogPost.GetPostsUseCase,
	getPostBySlug blogPost.GetPostBySlugUseCase,
) *BlogPostHandler {
	return &BlogPostHandler{
		createUseCase:       create,
		publishUseCase:      publish,
		updateUseCase:       update,
		deleteUseCase:       delete,
		getPostsUseCase:     getPosts,
		getPostBySlugUseCase: getPostBySlug,
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

	post, err := h.createUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
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
func currentUserID(c *gin.Context) string {
	if id, ok := c.Get("user_id"); ok {
		return id.(string)
	}
	return c.GetHeader("X-User-ID")
}

func (h *BlogPostHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido, deve ser um UUID válido"})
		return
	}

	userID := currentUserID(c)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "usuário não autenticado"})
		return
	}

	var input blogPost.UpdatePostInputDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	if err := h.updateUseCase.Execute(c.Request.Context(), id, userID, input); err != nil {
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

	userID := currentUserID(c)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "usuário não autenticado"})
		return
	}

	if err := h.deleteUseCase.Execute(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post removido com sucesso (soft delete)!"})
}

// GetPosts godoc
// @Summary      Listar posts publicados
// @Description  Retorna posts publicados com paginação
// @Tags         posts
// @Produce      json
// @Param        limit   query  int  false  "Limite (1-100, padrão 10)"
// @Param        offset  query  int  false  "Deslocamento (padrão 0)"
// @Success      200  {array}   blogPost.PostOutput
// @Failure      500  {object}  map[string]string
// @Router       /posts [get]
func (h *BlogPostHandler) GetPosts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit < 1 {
		limit = 1
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	posts, err := h.getPostsUseCase.Execute(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// GetPostBySlug godoc
// @Summary      Buscar post por slug
// @Description  Retorna um post publicado pelo slug
// @Tags         posts
// @Produce      json
// @Param        slug  path  string  true  "Slug do post"
// @Success      200  {object}  blogPost.PostOutput
// @Failure      404  {object}  map[string]string
// @Router       /posts/slug/{slug} [get]
func (h *BlogPostHandler) GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug é obrigatório"})
		return
	}

	post, err := h.getPostBySlugUseCase.Execute(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post não encontrado"})
		return
	}

	c.JSON(http.StatusOK, post)
}
