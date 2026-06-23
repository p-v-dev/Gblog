package http

import (
	"Gblog/internal/comment"
	"Gblog/internal/infra"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	createUseCase comment.CreateCommentUseCase
	listUseCase   comment.ListCommentsUseCase
	deleteUseCase comment.DeleteCommentUseCase
}

func NewCommentHandler(create comment.CreateCommentUseCase, list comment.ListCommentsUseCase, delete comment.DeleteCommentUseCase) *CommentHandler {
	return &CommentHandler{createUseCase: create, listUseCase: list, deleteUseCase: delete}
}

// Create godoc
// @Summary      Criar comentário
// @Description  Adiciona um comentário a um post
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        id           path  string                    true  "ID do post"
// @Param        comment      body  comment.CreateCommentInput true  "Conteúdo do comentário"
// @Param        Authorization header string false "Bearer {token}"
// @Success      201  {object}  comment.CommentOutput
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      422  {object}  map[string]string
// @Router       /posts/{id}/comments [post]
func (h *CommentHandler) Create(c *gin.Context) {
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do post é obrigatório"})
		return
	}

	userID := infra.CurrentUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
		return
	}

	var input comment.CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	input.PostID = postID
	input.UserID = userID

	output, err := h.createUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, output)
}

// List godoc
// @Summary      Listar comentários de um post
// @Description  Retorna todos os comentários ativos de um post
// @Tags         comments
// @Produce      json
// @Param        id  path  string  true  "ID do post"
// @Success      200  {array}  comment.CommentOutput
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /posts/{id}/comments [get]
func (h *CommentHandler) List(c *gin.Context) {
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do post é obrigatório"})
		return
	}

	comments, err := h.listUseCase.Execute(c.Request.Context(), postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// Delete godoc
// @Summary      Deletar comentário
// @Description  Remove o comentário (apenas o autor pode deletar)
// @Tags         comments
// @Produce      json
// @Param        id           path   string  true  "ID do comentário"
// @Param        Authorization header string false "Bearer {token}"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /comments/{id} [delete]
func (h *CommentHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do comentário é obrigatório"})
		return
	}

	userID := infra.CurrentUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
		return
	}

	if err := h.deleteUseCase.Execute(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comentário removido com sucesso"})
}
