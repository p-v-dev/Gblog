package http

import (
	"Gblog/internal/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	createUseCase user.CreateUserUseCase
	getUseCase    user.GetUserUseCase
	updateUseCase user.UpdateUserUseCase
	deleteUseCase user.DeleteUserUseCase
}

func NewUserHandler(create user.CreateUserUseCase, get user.GetUserUseCase, update user.UpdateUserUseCase, delete user.DeleteUserUseCase) *UserHandler {
	return &UserHandler{createUseCase: create, getUseCase: get, updateUseCase: update, deleteUseCase: delete}
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

// Get godoc
// @Summary      Buscar usuário por ID
// @Description  Retorna informações públicas de um usuário
// @Tags         users
// @Produce      json
// @Param        id  path  string  true  "ID do usuário (UUID)"
// @Success      200  {object}  user.UserOutput
// @Failure      404  {object}  map[string]string
// @Router       /users/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID é obrigatório"})
		return
	}

	output, err := h.getUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

// Update godoc
// @Summary      Atualizar perfil do usuário
// @Description  Atualiza nome e/ou email de um usuário ativo
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path  string                    true  "ID do usuário (UUID)"
// @Param        user  body  user.UpdateUserInput      true  "Dados para atualização"
// @Success      200   {object}  user.UserOutput
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Router       /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID é obrigatório"})
		return
	}

	var input user.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
		return
	}

	output, err := h.updateUseCase.Execute(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

// Delete godoc
// @Summary      Desativar conta de usuário
// @Description  Realiza a exclusão lógica (soft delete) de um usuário
// @Tags         users
// @Produce      json
// @Param        id  path  string  true  "ID do usuário (UUID)"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID é obrigatório"})
		return
	}

	if err := h.deleteUseCase.Execute(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuário desativado com sucesso"})
}
