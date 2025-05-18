package main

import (
	"github.com/stripe/stripe-go/v79"
	"gopkg.in/mail.v2"
	"optimaHurt/constAndVars"
	"optimaHurt/endpoints"
	"os"
)


// potem automatycznie dołączam usera do requesta
func main() {
	// podzielić to na mikroserwisy
	connection := endpoints.ConnectToDB(os.Getenv("CONNECTION_STRING"))
	defer func() {
		connection.Disconnect(constAndVars.ContextBackground)
	}()

	constAndVars.EmailDialer = mail.NewDialer("smtp.gmail.com", 587, "optimahurtcorp@gmail.com", os.Getenv("EMAIL_PASSWORD"))

	stripe.Key = os.Getenv("STRIPE_KEY")
	r := endpoints.MakeRouter()

	r.Run("0.0.0.0:" + os.Getenv("PORT"))
	return
}
