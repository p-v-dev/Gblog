package http

import (
	"Gblog/internal/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	createUseCase user.CreateUserUseCase
}

func NewUserHandler(create user.CreateUserUseCase) *UserHandler {
	return &UserHandler{createUseCase: create}
}

// Create godoc
// @Summary      Criar um usuário
// @Description  Cria um novo usuário no sistema
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      user.CreateUserInput  true  "Dados do usuário"
// @Success      201   {object}  user.UserOutput
// @Failure      400   {object}  map[string]string
// @Router       /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var input user.CreateUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
		return
	}

	output, err := h.createUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, output)
}
