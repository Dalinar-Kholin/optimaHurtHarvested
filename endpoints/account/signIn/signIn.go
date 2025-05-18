package signIn

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	. "optimaHurt/constAndVars"
	"optimaHurt/stringCheckers"
	"optimaHurt/user"
	"strconv"
	"strings"
)

// tworzy usera w bazie, pozwala na
func SignIn(c *gin.Context) {

	var reqData user.SignInBodyData
	if err := json.NewDecoder(c.Request.Body).Decode(&reqData); err != nil {
		c.JSON(200, gin.H{
			"error": "złe dane",
		})
		return
	}

	sendErrorResponse := func(err error) {
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
	}

	trimmer := func(s string) string {
		return strings.Trim(s, " \t")
	}
	checkers := []struct {
		toValidate string
		validator  func(data string) error
	}{ // walidacja danych
		{
			reqData.Password,
			stringCheckers.CheckPassword,
		},
		{
			reqData.Nip,
			stringCheckers.CheckNip,
		},
		{
			reqData.Email,
			stringCheckers.CheckEmail,
		},
		{
			reqData.Username,
			stringCheckers.CheckUsername,
		},
		{
			reqData.CompanyName,
			stringCheckers.CheckCompanyName,
		},
		{
			reqData.Street,
			func(data string) error {
				if strings.ContainsAny(data, "!@#$%^&*()+=[]{}|\\:;\"'<>,.?/~`") {
					return errors.New("nr  zwiera niepoprawne znaki")
				}
				return nil
			},
		},
		{
			reqData.Nr,
			func(data string) error {
				_, err := strconv.Atoi(data)
				if err != nil {
					return errors.New("niepoprawne znaki w nr domu")
				}

				return nil
			},
		},
	}

	for _, x := range checkers {
		if err := x.validator(trimmer(x.toValidate)); err != nil {
			sendErrorResponse(err)
			return
		}
	} // walidacja danych

	connection := DbConnect.Collection(UserCollection)
	var userInDb user.DataBaseUserObject
	if connection.FindOne(ContextBackground, bson.M{
		"$or": []bson.M{
			{"nip": reqData.Nip},
			{"username": reqData.Username},
			{"email": reqData.Email}, // email używany jest do rozliczania więc całkiem ważne
		},
	}).Decode(&userInDb) == nil {
		if userInDb.Username == reqData.Username {
			c.JSON(200, gin.H{
				"error": "nazwa użytkownika aktualnie istnieje",
			})
			return
		} else {
			fmt.Printf("user in Db := %v", userInDb)
			c.JSON(200, gin.H{
				"error": "dany nip już występuję",
			})
			return
		}
	}

	newUser := user.DataBaseUserObject{
		Id:       primitive.NewObjectID(),
		Email:    reqData.Email,
		Username: reqData.Username,
		Password: reqData.Password,
		CompanyData: user.CompanyData{
			Name: reqData.CompanyName,
			Nip:  reqData.Nip,
			Address: user.Address{
				Street: reqData.Street,
				Nr:     reqData.Nr,
			},
		},
		AccountStatus: user.Active,
	}

	if _, err := connection.InsertOne(ContextBackground, newUser); err != nil {
		c.JSON(500, gin.H{
			"error": "błąd servera",
		})
		return
	}
	//ICOM: Na razie każdy może korzystać ile chce bez problemu

	// message := user.UserMessage{UserId: newUser.Id, Message: "aby odblokować możliwość sprawdzania przejdź do płatności i rozpocznij okres próbny"}

	//_, _ = DbConnect.Collection(UserMessageCollection).InsertOne(ContextBackground, message)

	c.JSON(200, gin.H{
		"result": "dodano użytkownika",
	})

}
