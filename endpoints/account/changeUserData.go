package account

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	. "optimaHurt/constAndVars"
	"optimaHurt/stringCheckers"
	"optimaHurt/user"
)

type DataToInsert struct {
	NewHurtDataTab []user.UserCreds `json:"newHurtData"`
	NewAccountData string           `json:"newAccountData"`
}

func handleHurtData(userInstance *user.DataBaseUserObject, data []user.UserCreds) {
	for _, i := range data {
		isIn := false
		for j, x := range userInstance.Creds {
			if x.HurtName == i.HurtName {
				userInstance.Creds[j] = i
				isIn = true
			}
		}
		if !isIn {
			userInstance.Creds = append(userInstance.Creds, i)
		}
	}
}

func handleAccountData(userInstance *user.DataBaseUserObject, data string) {
	userInstance.Password = data
}

func ChangeUserData(c *gin.Context) {
	// zanim się połączymy z bazą sprawdzmy czy to ma sens
	var data DataToInsert
	reader := json.NewDecoder(c.Request.Body)
	err := reader.Decode(&data)
	if err != nil {
		c.JSON(200, gin.H{
			"error": "złe dane",
		})
		return
	}

	auth := c.Request.Header.Get("Authorization") // w middlewear sprawdziliśmy i wiemy że jest poprawnie zdefiniowne
	userInstance := Users[auth]

	con := DbConnect.Collection(UserCollection)
	session, err := DbClient.StartSession() // rozpoczęci sesji aby potem zacząć transakcję
	if err != nil {
		c.JSON(501, gin.H{"error": "server stupido"})
		return
	}
	defer session.EndSession(ContextBackground)

	err = session.StartTransaction()
	if err != nil {
		c.JSON(501, gin.H{"error": "server stupido"})
		return
	}
	err = mongo.WithSession(ContextBackground, session, func(sc mongo.SessionContext) error {

		var dataBaseUser user.DataBaseUserObject
		err = con.FindOne(ContextBackground, bson.M{
			"_id": userInstance.Id,
		}).Decode(&dataBaseUser)
		if err != nil {
			return err
		}

		for _, x := range data.NewHurtDataTab {
			if err := stringCheckers.CheckUsername(x.Login); err != nil {
				if err := stringCheckers.CheckEmail(x.Login); err != nil {
					return err
				}
			}
			if err := stringCheckers.CheckPassword(x.Password); err != nil {
				return err
			}
		}
		if err := stringCheckers.CheckPassword(data.NewAccountData); data.NewAccountData != "" && err != nil {
			return err
		}

		if data.NewHurtDataTab != nil {
			handleHurtData(&dataBaseUser, data.NewHurtDataTab)
		}

		if data.NewAccountData != "" {
			handleAccountData(&dataBaseUser, data.NewAccountData)
		}

		var id user.DataBaseUserObject
		err = con.FindOneAndReplace(ContextBackground, bson.M{
			"_id": userInstance.Id,
		}, dataBaseUser).Decode(&id)
		if err != nil {
			return errors.New("nie znaleziono użytkownika")
		}
		return nil
	})

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := session.CommitTransaction(ContextBackground); err != nil {
		c.JSON(500, gin.H{
			"error": "server error",
		})
		return
	}
	delete(Users, auth)
	c.JSON(200, gin.H{
		"result": "udało się dodać hasło",
	})
	return
}
