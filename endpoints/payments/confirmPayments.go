package payments

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v79"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	. "optimaHurt/constAndVars"
	"optimaHurt/user"
)

func ConfirmPayment(c *gin.Context) {
	var event stripe.Event
	if err := json.NewDecoder(c.Request.Body).Decode(&event); err != nil {
		c.Status(400)
		return
	}

	switch event.Type {
	case "customer.subscription.deleted": // kiedy to się dzieje anulujemy użytkownikowi usługę
		var delInfo stripe.Subscription

		if err := json.Unmarshal(event.Data.Raw, &delInfo); err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}
		var stripeInfo user.StripeUserInfo
		if err := DbConnect.Collection(StripeCollection).FindOneAndDelete(ContextBackground, bson.M{
			"subscriptionId": delInfo.ID,
		}).Decode(&stripeInfo); err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}
		if err := DbConnect.
			Collection(UserCollection).
			FindOneAndUpdate(ContextBackground,
				bson.M{"_id": stripeInfo.UserId},
				bson.M{"$set": bson.M{"accountStatus": user.Inactive}}).Err(); err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "all is GIT",
		})
		return
	case "checkout.session.completed":
		var session stripe.CheckoutSessionParams
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		idString := session.Metadata["userId"]
		id, err := primitive.ObjectIDFromHex(idString)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		conn := DbConnect.Collection(UserCollection)
		if _, err = conn.UpdateOne(ContextBackground, bson.M{"_id": id}, bson.M{"$set": bson.M{"accountStatus": user.Active}}); err != nil {
			c.JSON(400, gin.H{
				"error": err,
			})
			return
		}
		subscriptionID := event.Data.Object["subscription"].(string)
		info := user.StripeUserInfo{
			UserId:         id,
			SubscriptionId: subscriptionID,
		}
		_, _ = DbConnect.Collection(StripeCollection).InsertOne(ContextBackground, info)
		messageConn := DbConnect.Collection(UserMessageCollection)
		message := user.UserMessage{UserId: id, Message: "płatność się udała"}
		messageConn.InsertOne(ContextBackground, message)
		return
	}
}
