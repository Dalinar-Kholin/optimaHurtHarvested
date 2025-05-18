package forgotPassword

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mail.v2"
	"net/http"
	. "optimaHurt/constAndVars"
	"optimaHurt/stringCheckers"
	"optimaHurt/user"
	"os"
)

type ForgotPasswordInDb struct {
	UserId primitive.ObjectID `bson:"userId"`
	Token  string             `bson:"token"`
}

func ForgotPassword(c *gin.Context) {

	email := c.Request.URL.Query().Get("email")

	if err := stringCheckers.CheckEmail(email); err != nil {
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
		return
	}
	var userInDb user.DataBaseUserObject

	if err := DbConnect.Collection(UserCollection).FindOne(ContextBackground, bson.M{
		"email": email,
	}).Decode(&userInDb); err != nil {
		c.JSON(200, gin.H{
			"error": "zły adres mail",
		})
		return
	}

	token := make([]byte, 64)
	if _, err := rand.Read(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cant generate token"})
		return
	}
	m := mail.NewMessage()
	m.SetHeader("To", email)
	m.SetHeader("From", "optimaHurtCorp@gmail.com")
	m.SetHeader("Subject", "zapomniane hasło")

	shielded := base64.URLEncoding.EncodeToString(token)

	m.SetBody("text/plain; charset=UTF-8", fmt.Sprintf("aby zrestartowac haslo prosze wejsc pod strone https://optimahurt-hayvfpjoza-lm.a.run.app/resetPassword?token=%s\n", shielded))

	forgot := ForgotPasswordInDb{UserId: userInDb.Id, Token: shielded}

	if _, err := DbConnect.Collection(ResetPassword).InsertOne(ContextBackground, forgot); err != nil {
		c.JSON(500, gin.H{
			"error": "server error",
		})
	}

	if err := mail.NewDialer("smtp.gmail.com", 587, "optimahurtcorp@gmail.com", os.Getenv("EMAIL_PASSWORD")).DialAndSend(m); err != nil {
		fmt.Printf("%v", err)
		c.JSON(500, gin.H{
			"error": "nie udało się wysłać maila",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "sprawdź email",
	})
}
