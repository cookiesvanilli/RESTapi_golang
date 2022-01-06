package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

// прослойка, которая парсит токены из запроса и предоставляет доступ к группе endpoints /api
func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header") // 401 - пользователь не авторизирован
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 { // массив длиною в 2 элемента
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}
	// parse token
	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(userCtx, userId)
}
