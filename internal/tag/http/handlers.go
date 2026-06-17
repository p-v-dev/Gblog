package http

import (
	"Gblog/internal/tag"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	createUseCase tag.CreateTagUseCase
	listUseCase   tag.ListTagsUseCase
}

func NewTagHandler(create tag.CreateTagUseCase, list tag.ListTagsUseCase) *TagHandler {
	return &TagHandler{createUseCase: create, listUseCase: list}
}

// Create godoc
// @Summary      Criar uma tag
// @Description  Cria uma nova tag ou retorna a existente se o nome já existir
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        tag  body  tag.CreateTagInput  true  "Nome da tag"
// @Success      201  {object}  tag.TagOutput
// @Failure      400  {object}  map[string]string
// @Router       /tags [post]
func (h *TagHandler) Create(c *gin.Context) {
	var input tag.CreateTagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	output, err := h.createUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, output)
}

// List godoc
// @Summary      Listar tags
// @Description  Retorna todas as tags cadastradas
// @Tags         tags
// @Produce      json
// @Success      200  {array}  tag.TagOutput
// @Router       /tags [get]
func (h *TagHandler) List(c *gin.Context) {
	tags, err := h.listUseCase.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tags)
}
