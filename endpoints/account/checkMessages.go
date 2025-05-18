package account

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	. "optimaHurt/constAndVars"
	"optimaHurt/user"
)

func CheckMessages(c *gin.Context) {

	userInstance := Users[c.Request.Header.Get("Authorization")]

	conn := DbConnect.Collection(UserMessageCollection)

	var message user.UserMessage

	if err := conn.FindOneAndDelete(ContextBackground, bson.M{
		"userId": userInstance.Id,
	}).Decode(&message); err != nil {
		c.JSON(200, gin.H{
			"message": "",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": message.Message,
	})

}
