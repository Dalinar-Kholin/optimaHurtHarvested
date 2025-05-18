package account

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"optimaHurt/constAndVars"
)

func TestCookie(c *gin.Context) {
	Response := func(c *gin.Context, response bool, status int) {
		c.JSON(status, gin.H{
			"response": response,
		})
	}
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		Response(c, false, http.StatusUnauthorized)
		return
	}
	_, ok := constAndVars.Users[token]
	if !ok {
		Response(c, false, http.StatusUnauthorized)
		return
	}
	Response(c, true, http.StatusOK) // ciasteczko jest prawidłowe
}
