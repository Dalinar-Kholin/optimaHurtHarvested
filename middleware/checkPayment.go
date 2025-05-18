package middleware

import (
	"github.com/gin-gonic/gin"
	. "optimaHurt/constAndVars"
	"optimaHurt/user"
)

func CheckPayment(c *gin.Context) {

	// za wolne, bardzo spowalnia zapytanie, zamiast tego dodać pole podczas logowania
	token := c.Request.Header.Get("Authorization") // wiem że będzie bo jesteśmy już za innymi bramkami
	userInstance := Users[token]

	if userInstance.AccountStatus != user.Active {
		c.JSON(400, gin.H{
			"error": "make Payment",
		})
		c.Abort()
		return
	}
}
