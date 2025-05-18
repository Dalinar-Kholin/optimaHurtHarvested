package account

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"net/http"
	. "optimaHurt/constAndVars"
	"optimaHurt/hurtownie"
	"optimaHurt/hurtownie/factory"
	"optimaHurt/stringCheckers"
	"optimaHurt/user"
	"sync"
)

func makeNewUser(dbUser user.DataBaseUserObject) (userInstance *user.User, availableHurts hurtownie.HurtName, resultsLogin []ChannelResponse) {
	var hurtTab []hurtownie.IHurt
	var validHurt []hurtownie.IHurt
	var wg sync.WaitGroup

	/*proxyURL, err := url.Parse("http://127.0.0.1:8000")
	if err != nil {
		fmt.Println("Błąd parsowania URL proxy:", err)
		return
	}

	// Konfiguracja Transport z Proxy
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}*/

	client := &http.Client{
		//Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	ch := make(chan ChannelResponse)
	var credsEnum hurtownie.HurtName = 0
	for _, creds := range dbUser.Creds {
		instance, _ := factory.HurtFactory(creds.HurtName)
		hurtTab = append(hurtTab, instance)
		credsEnum += creds.HurtName
		wg.Add(1)
		go func(hurt *hurtownie.IHurt, wg *sync.WaitGroup, ch chan<- ChannelResponse) {
			defer wg.Done()
			res := (*hurt).TakeToken(creds.Login, creds.Password, client)
			name := (*hurt).GetName()
			if res {
				validHurt = append(validHurt, *hurt)
				ch <- ChannelResponse{
					Hurt:    name,
					Success: true,
				}
			} else {
				ch <- ChannelResponse{
					Hurt:    name,
					Success: false,
				}
			}
		}(&instance, &wg, ch)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	res := make([]ChannelResponse, len(hurtTab))
	valid := make([]ChannelResponse, 0)
	i := 0
	for x := range ch {
		if x.Success {
			availableHurts += x.Hurt
			valid = append(valid, x)
		}
		res[i] = x
		i++
	}
	userInstance = &user.User{Client: client, Hurts: validHurt, Creds: dbUser.Creds}
	resultsLogin = res
	return
}

type ChannelResponse struct {
	Hurt    hurtownie.HurtName `json:"hurt"`
	Success bool               `json:"success"`
}

func Login(c *gin.Context) {
	connection := DbConnect.Collection(UserCollection)

	request, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant read body"})
		return
	}
	defer c.Request.Body.Close()

	var reqBody user.LoginBodyData

	if err = json.Unmarshal(request, &reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}

	if err := stringCheckers.CheckUsername(reqBody.Username); err != nil {
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Printf("hasło := %v\n", reqBody.Password)
	if err := stringCheckers.CheckPassword(reqBody.Password); err != nil {
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
		return
	}

	var dataBaseResponse user.DataBaseUserObject

	if err = connection.FindOne(ContextBackground, bson.M{
		"username": reqBody.Username,
		"password": reqBody.Password,
	}).Decode(&dataBaseResponse); err != nil {
		c.JSON(200, gin.H{"error": "nieprawidłowe dane"})
		return
	}

	userInstance, loggedHurts, loginLog := makeNewUser(dataBaseResponse)
	userInstance.Id = dataBaseResponse.Id
	userInstance.AccountStatus = dataBaseResponse.AccountStatus
	newToken := make([]byte, 64)
	if _, err := rand.Read(newToken); err != nil {
		c.JSON(500, gin.H{"error": "cant generate token"})
		return
	}
	shieldedToken := base64.URLEncoding.EncodeToString(newToken)
	//shieldedToken := base64.StdEncoding.EncodeToString(newToken) // token jest użytkowany tylko podczas sesji, więc nie ma potrzeby przechowywania go w bazie danych
	Users[shieldedToken] = userInstance

	loginResult := LoginResult{
		Result:         loginLog,
		Token:          shieldedToken,
		AvailableHurts: loggedHurts,
		Status:         dataBaseResponse.AccountStatus,
		CompanyName:    dataBaseResponse.CompanyData.Name,
	}

	c.JSON(200, loginResult)
}

type LoginResult struct {
	Result         []ChannelResponse  `json:"result"`
	Token          string             `json:"token"`
	AvailableHurts hurtownie.HurtName `json:"availableHurts"`
	Status         user.AccountStatus `json:"accountStatus"`
	CompanyName    string             `json:"companyName"`
}
