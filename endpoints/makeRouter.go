package endpoints

import (
	"github.com/gin-gonic/gin"
	"optimaHurt/constAndVars"
	"optimaHurt/endpoints/account"
	"optimaHurt/endpoints/account/forgotPassword"
	"optimaHurt/endpoints/account/signIn"
	"optimaHurt/endpoints/orders"
	"optimaHurt/endpoints/payments"
	"optimaHurt/endpoints/takePrices"
	"optimaHurt/middleware"
)

func MakeRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.AddHeaders)
	r.Static("/assets", "./frontend/dist/assets")

	// Obsługa głównego pliku index.html
	r.StaticFile("/", "./frontend/dist/index.html")

	// Obsługa aplikacji typu SPA - przekierowanie wszystkich nieznalezionych ścieżek do index.html
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	// Dodanie trasy API
	api := r.Group("/api")
	{
		api.GET("isAlive", func(c *gin.Context) {
			c.JSON(200, true)
		})

		api.POST("/checkCredentials", account.CheckCredentials)

		api.POST("/checkCookie", account.TestCookie)

		api.GET("/logout", func(c *gin.Context) {
			auth := c.Request.Header.Get("Authorization")
			if auth == "" {
				return
			}
			delete(constAndVars.Users, auth)
		})

		api.POST("/takePrices", middleware.CheckToken, middleware.CheckHurtTokenCurrency, middleware.CheckPayment, takePrices.TakeMultiple) // get nie może mieć body, więc robimy post
		api.GET("/takePrice", middleware.CheckToken, middleware.CheckHurtTokenCurrency, middleware.CheckPayment, takePrices.TakePrice)
		api.POST("/makeOrder", middleware.CheckToken, middleware.CheckHurtTokenCurrency, middleware.CheckPayment, orders.MakeOrder)

		api.GET("/forgotPassword", forgotPassword.ForgotPassword)
		api.POST("/resetPassword", forgotPassword.ResetPasswordFunc)

		api.POST("/login", account.Login)
		api.POST("/signIn", signIn.SignIn)

		api.GET("/messages", middleware.CheckToken, account.CheckMessages)

		api.POST("/payment/stripe", middleware.CheckToken, payments.MakePayment)
		api.POST("/payment/stripe/webhook/confirm", payments.ConfirmPayment)
		api.GET("/payment/stripe/cancel", middleware.CheckToken, payments.EndSubscription)

		api.PATCH("/changeUserData", middleware.CheckToken, account.ChangeUserData)
	}
	return r
}
