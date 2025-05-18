package payments

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/subscription"
	"go.mongodb.org/mongo-driver/bson"
	. "optimaHurt/constAndVars"
	"optimaHurt/user"
	"time"
)

func EndSubscription(c *gin.Context) {

	token := c.Request.Header.Get("Authorization")
	userInstance := Users[token]

	var stripeInfo user.StripeUserInfo

	if err := DbConnect.Collection(StripeCollection).FindOne(ContextBackground, bson.M{
		"userId": userInstance.Id,
	}).Decode(&stripeInfo); err != nil {
		fmt.Printf("error przy anulowaniu subskrypcji := %v", err.Error())
		c.JSON(500, gin.H{
			"error": "server stupido",
		})
		return
	}

	params := &stripe.SubscriptionParams{CancelAtPeriodEnd: stripe.Bool(true)}
	result, err := subscription.Update(stripeInfo.SubscriptionId, params)
	if err != nil {
		fmt.Printf("error przy anulowaniu subskrypcji := %v\n", err.Error())
		c.JSON(500, gin.H{
			"error": "cant update subscription",
		})
		return
	}

	date := time.Unix(result.CurrentPeriodEnd, 0)

	// Wyświetlanie daty
	c.JSON(200, gin.H{
		"messaage": fmt.Sprintf("konto ma zapewniony dostęp do %s", date),
	})

}
