package account

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"optimaHurt/hurtownie"
	"optimaHurt/hurtownie/factory"
)

type DataToCheck struct {
	Username string             `json:"username"`
	Password string             `json:"password"`
	HurtName hurtownie.HurtName `json:"hurtName"`
}

func CheckCredentials(c *gin.Context) {
	var Creds DataToCheck
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&Creds)
	if err != nil {
		c.JSON(200, gin.H{
			"error": "bad Request Body",
		})
		return
	}
	fmt.Printf("dane lowowania do hurtowni := %v\n", Creds)
	hurt, err := factory.HurtFactory(Creds.HurtName)
	if err != nil {
		c.JSON(200, gin.H{
			"error": "bad Hurt enum",
		})
		return
	}
	res := hurt.TakeToken(Creds.Username, Creds.Password, &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	})

	if !res {
		c.JSON(200, gin.H{
			"error": "nie znaleziono konta dla podanych danych",
		})
		return
	}
	/*ICOM:  wiemy że dane sa prawidłowe, teraz na frontendzie dodamy je do poprawnych i użytkownik przy zapisywaniu zmian prześle wszystkie poprawne dane ograniczamy tym ilość zapytań do bazy danych które są wolne i drogie*/
	c.JSON(200, gin.H{
		"result": "udao się dodać hasło",
	})

}
