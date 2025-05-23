package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"optimaHurt/constAndVars"
	"optimaHurt/hurtownie"
	"optimaHurt/user"
	"sync"
)

func CheckHurtTokenCurrency(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		c.JSON(400, gin.H{
			"error": "where Token?",
		})
		return
	}
	var ok bool
	var userInstance *user.User
	if userInstance, ok = constAndVars.Users[token]; !ok {
		c.JSON(400, gin.H{
			"error": "where logowanie?",
		})
	}

	var wg sync.WaitGroup

	for _, hurt := range userInstance.Hurts {
		wg.Add(1)
		go func(wg *sync.WaitGroup, hurt hurtownie.IHurt) {
			defer wg.Done()
			if !hurt.CheckToken(userInstance.Client) {
				fmt.Printf("refresh tokena hurtowni := %v", hurt.GetName())
				if !hurt.RefreshToken(userInstance.Client) {
					userCred := userInstance.TakeHurtCreds(hurt.GetName())
					if !hurt.TakeToken(userCred.Login, userCred.Password, userInstance.Client) {
						c.JSON(400, gin.H{
							"error": "nie udało się zalogować do hurtowni",
						})
					}
				}
			}
		}(&wg, hurt)
	}
	wg.Wait()
}
