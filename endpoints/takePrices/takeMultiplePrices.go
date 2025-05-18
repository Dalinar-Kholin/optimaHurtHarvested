package takePrices

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"optimaHurt/constAndVars"
	"optimaHurt/hurtownie"
	"optimaHurt/stringCheckers"
	"sync"
)

func TakeMultiple(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	userInstance := constAndVars.Users[token]

	var list hurtownie.WishList

	err := json.NewDecoder(c.Request.Body).Decode(&list)
	defer c.Request.Body.Close()

	if err != nil {
		c.JSON(200, gin.H{
			"error": "złe dane",
		})
		return
	}

	for _, x := range list.Items {
		if err := stringCheckers.CheckEan(x.Ean); err != nil {
			c.JSON(200, gin.H{
				"error": fmt.Sprintf("w ean := %v wystąpił błąd := %v", x.Ean, err.Error()),
			})
			return
		}
	}

	var wg sync.WaitGroup
	ch := make(chan Result)

	for _, hurt := range userInstance.Hurts {
		wg.Add(1)
		go func(hurt *hurtownie.IHurt, wg *sync.WaitGroup, ch chan<- Result) {
			defer wg.Done()
			res, err := (*hurt).SearchMany(list, userInstance.Client)
			if err != nil {
				ch <- Result{
					HurtName: (*hurt).GetName(),
					Result:   nil,
				}
				return
			}
			ch <- Result{
				HurtName: (*hurt).GetName(),
				Result:   res}
		}(&hurt, &wg, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()
	result := make([]interface{}, len(userInstance.Hurts))
	i := 0

	for x := range ch {
		result[i] = x
		i++
	}
	c.JSON(200, result)
}

type Result struct {
	HurtName hurtownie.HurtName
	Result   []hurtownie.SearchManyProducts
}
