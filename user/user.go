package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"optimaHurt/hurtownie"
)

type AccountStatus int

const (
	Inactive AccountStatus = iota // użykownik który anulował subskrypcje
	New                           // przysyługuje mu zniżna na pierwsze zamówienie, przy anulowaniu subskrypcji i rozpoczęciu nowej, nie ma już 2 tygodni na sprawdzenie
	Active                        // ładnie płacący user
)

type SignInBodyData struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	CompanyName string `json:"companyName"`
	Nip         string `json:"nip"`
	Street      string `json:"street"`
	Nr          string `json:"nr"`
}

type LoginBodyData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserMessage struct {
	UserId  primitive.ObjectID `bson:"userId"`
	Message string             `bson:"message"`
}

type User struct {
	Id            primitive.ObjectID
	Client        *http.Client
	Hurts         []hurtownie.IHurt
	Creds         []UserCreds
	AccountStatus AccountStatus `bson:"subscriptionTier" json:"subscriptionTier"`
}

func (u *User) TakeHurtCreds(name hurtownie.HurtName) UserCreds {
	for _, i := range u.Creds {
		if i.HurtName == name {
			return i
		}
	}
	return UserCreds{}
}

type StripeUserInfo struct {
	UserId         primitive.ObjectID `bson:"userId"`
	SubscriptionId string             `bson:"subscriptionId"`
}

type DataBaseUserObject struct {
	Id            primitive.ObjectID `bson:"_id" json:"_id"`
	Email         string             `bson:"email" json:"email"`
	Username      string             `bson:"username" json:"username"`
	Password      string             `bson:"password" json:"password"`
	CompanyData   CompanyData        `bson:"companyData" json:"companyData"`
	Creds         []UserCreds        `bson:"creds" json:"creds"`
	AccountStatus AccountStatus      `bson:"accountStatus" json:"accountStatus"`
}

type UserCreds struct {
	HurtName hurtownie.HurtName `json:"hurtName" bson:"hurtName"`
	Login    string             `json:"login" bson:"login"`
	Password string             `json:"password" bson:"password"`
}

type Address struct {
	Street string `bson:"street"  bson:"street"`
	Nr     string `bson:"nr" json:"nr"`
}

type CompanyData struct {
	Name    string  `bson:"CompanyName" json:"companyName"`
	Nip     string  `bson:"nip" json:"nip"`
	Address Address `bson:"address" json:"address"`
}
