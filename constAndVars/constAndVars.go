package constAndVars

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mail.v2"
	"optimaHurt/user"
)

const (
	DbName                = "optiHurt"
	ResetPassword         = "resetPasswordRequest"
	UserCollection        = "users"
	UserMessageCollection = "UserMessage"
	StripeCollection      = "stripeInfo"
)

var (
	EmailDialer       *mail.Dialer
	Users             map[string]*user.User = make(map[string]*user.User) // mapuje id na usera -- zakładam że userów nie będzie jakoś strasznie dużo
	ContextBackground context.Context       = context.TODO()
	DbConnect         *mongo.Database
	DbClient          *mongo.Client
)

func ExportedFunction() {
	// Możesz tu umieścić jakikolwiek kod, nawet zostawić funkcję pustą
}
