package takePrices

import (
	"github.com/gin-gonic/gin"
	"optimaHurt/constAndVars"
	"optimaHurt/hurtownie"
	"optimaHurt/stringCheckers"
	"optimaHurt/user"
	"sync"
)

func TakePrice(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	userInstance := constAndVars.Users[token]

	ean := c.Query("ean")

	if err := stringCheckers.CheckEan(ean); err != nil {
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
		return
	}

	var wg sync.WaitGroup
	ch := make(chan SearchResult)
	for _, hurt := range userInstance.Hurts {
		wg.Add(1)
		go func(hurt *hurtownie.IHurt, wg *sync.WaitGroup, ch chan<- SearchResult) {
			defer wg.Done()
			res, err := (*hurt).SearchProduct(ean, userInstance.Client)
			if err != nil && err.Error() == "tokenError" {
				var creds user.UserCreds
				for _, i := range userInstance.Creds {
					if i.HurtName == (*hurt).GetName() {
						creds = i
					}
				}
				(*hurt).TakeToken(creds.Login, creds.Password, userInstance.Client) // tutaj zamienić na refreshToken, nie ma potrzeby ponownie pobierać tokena
				res, err = (*hurt).SearchProduct(ean, userInstance.Client)
			}
			if err != nil {
				ch <- SearchResult{Ean: ean, Result: nil, HurtName: (*hurt).GetName()}
				return
			}
			ch <- SearchResult{Ean: ean, Result: res, HurtName: (*hurt).GetName()}
		}(&hurt, &wg, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()
	result := make([]SearchResult, len(userInstance.Hurts))
	i := 0

	for x := range ch {
		result[i] = x
		i++
	}
	c.JSON(200, result)

}

type SearchResult struct {
	Ean      string             `json:"ean"`
	HurtName hurtownie.HurtName `json:"hurtName"`
	Result   interface{}        `json:"result"`
}
