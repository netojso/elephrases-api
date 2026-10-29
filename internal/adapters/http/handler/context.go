package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/netojso/elephrases-api/pkg"
)

// requestContext deriva do request um context com o timeout configurado.
func requestContext(c *gin.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		return context.WithCancel(c.Request.Context())
	}
	return context.WithTimeout(c.Request.Context(), timeout)
}

// errorStatus responde 404 para pkg.ErrNotFound e 500 para os demais erros.
func errorStatus(err error) int {
	if errors.Is(err, pkg.ErrNotFound) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}
