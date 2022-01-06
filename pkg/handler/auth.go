package handler

import (
	todo "github.com/cookiesvanilli/go_app"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var input todo.User

	if err := c.BindJSON(&input); err != nil { // BindJSON принимает ссылку на объект в который хотим распарсить тело json
		newErrorResponse(c, http.StatusBadRequest, err.Error()) //400 пользователь ввел некорректные данные
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()) //500 внутренняя ошибка на сервере
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil { // BindJSON принимает ссылку на объект в который хотим распарсить тело json
		newErrorResponse(c, http.StatusBadRequest, err.Error()) //400 пользователь ввел некорректные данные
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()) //500 внутренняя ошибка на сервере
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
