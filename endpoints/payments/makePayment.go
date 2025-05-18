package payments

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/checkout/session"
	"github.com/stripe/stripe-go/v79/customer"
	"go.mongodb.org/mongo-driver/bson"
	. "optimaHurt/constAndVars"
	"optimaHurt/user"
	"strconv"
)

const (
	monthlyDefault = iota
	yearly
)

const (
	renewalSubscription        = "price_1PonIR03bfZgIVzM5PIbJmee"
	newSubscription            = "price_1PonIR03bfZgIVzM5PIbJmee"
	yearlySubscription         = "price_1PonJq03bfZgIVzMYyhkT9ct"
	discountYearlySubscription = "ZYUEZYw0"
)

func MakePayment(c *gin.Context) {

	auth := c.Request.Header.Get("Authorization") // jesteśmy za bramką
	userInstance := Users[auth]

	prodNameString := c.Request.URL.Query().Get("prodName")
	if prodNameString == "" || len(prodNameString) > 3 {
		c.JSON(200, gin.H{
			"error": "bad request",
		})
		return
	}

	prodName, err := strconv.Atoi(prodNameString)
	if err != nil {
		fmt.Printf("%v\n", prodName)
		c.JSON(200, gin.H{
			"error": "bad request",
		})
		return
	}

	conn := DbConnect.Collection(UserCollection)

	var userInDb user.DataBaseUserObject

	if err := conn.FindOne(ContextBackground, bson.M{"_id": userInstance.Id}).Decode(&userInDb); err != nil {
		c.JSON(500, gin.H{
			"error": "server error",
		})
		return
	}

	customerParams := &stripe.CustomerParams{
		Email: stripe.String(userInDb.Email),
	}

	newCustomer, err := customer.New(customerParams)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "server error",
		})
		return
	}

	priceId := newSubscription
	if prodName == monthlyDefault {
		if userInDb.AccountStatus != user.New {
			priceId = renewalSubscription
		}
	} else if prodName == yearly {
		priceId = yearlySubscription
	}

	// Tworzenie sesji płatności Stripe
	params := &stripe.CheckoutSessionParams{
		Customer:           &newCustomer.ID,
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				Price:    stripe.String(priceId),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String("subscription"),
		SuccessURL: stripe.String("https://optimahurt-hayvfpjoza-lm.a.run.app/login"),
		CancelURL:  stripe.String("https://optimahurt-hayvfpjoza-lm.a.run.app/failed"),
		Metadata: map[string]string{
			"userId": Users[auth].Id.Hex(), // Przekazanie identyfikatora użytkownika
		},
	}

	if prodName == monthlyDefault && userInDb.AccountStatus == user.New {
		params.SubscriptionData = &stripe.CheckoutSessionSubscriptionDataParams{
			TrialPeriodDays: stripe.Int64(14)}
	} // 14 dni próbne

	if prodName == yearly {
		params.Discounts = []*stripe.CheckoutSessionDiscountParams{
			&stripe.CheckoutSessionDiscountParams{
				Coupon: stripe.String(discountYearlySubscription),
			},
		}
	}

	s, err := session.New(params)
	if err != nil {
		fmt.Printf("error %v", err)
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Zwrócenie ID sesji płatności
	c.JSON(200, gin.H{
		"id": s.ID,
	})
}
