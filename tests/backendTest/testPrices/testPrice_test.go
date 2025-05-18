package testPrices

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"optimaHurt/constAndVars"
	"optimaHurt/endpoints"
	"optimaHurt/endpoints/account"
	"optimaHurt/endpoints/takePrices"
	"optimaHurt/hurtownie"
	"optimaHurt/hurtownie/eurocash"
	"optimaHurt/hurtownie/tedi"
	"optimaHurt/tests/backendTest"
	"optimaHurt/user"
	"os"
	"testing"
)

func Test(t *testing.T) {
	// Tworzenie nowego routera Gin w trybie testowym

	router := endpoints.MakeRouter()
	connection := endpoints.ConnectToDB(os.Getenv("CONNECTION_STRING"))
	defer func() {
		connection.Disconnect(constAndVars.ContextBackground)
	}()

	tediEan := "5900541000000"     // kod ean żywca zdrój -- powinien być pernamentnie dostępny
	sotEan := "5900951311321"      // mars
	specjalEan := "5900541000000"  // kod ean żywca zdrój
	eurocashEan := "5900541000000" // kod ean żywca zdrój

	// ICOM: tworzenie danych do testów
	toTest := []struct {
		ean  string
		hurt hurtownie.HurtName
	}{
		{
			eurocashEan,
			hurtownie.Eurocash,
		},
		{
			specjalEan,
			hurtownie.Specjal,
		},
		{
			sotEan,
			hurtownie.Sot,
		},
		{
			tediEan,
			hurtownie.Tedi,
		},
	}

	for _, data := range toTest {
		if isOk, _ := testHurt(data.ean, "", router, t, data.hurt, false, false); isOk {
			t.Logf("unlogged user can see results")
		}
	}

	loginResult, err := login(backendTest.TestUserUsername, backendTest.TestUserPassword, router)
	if err != nil {
		t.Logf("%v", err)
		return
	}
	token := loginResult.Token
	for _, data := range toTest { // checking products that exist
		if isOk, _ := testHurt(data.ean, token, router, t, data.hurt, true, true); !isOk {
			t.Logf("user is logged and still dont work")
		}
	}
	for _, data := range toTest { // checking products that don't exist
		if isOk, _ := testHurt("40084015546738", token, router, t, data.hurt, true, false); isOk {
			t.Logf("product should not exist")
		}
	}
}

func login(username, password string, router *gin.Engine) (*account.LoginResult, error) {
	data := user.LoginBodyData{
		Username: username,
		Password: password,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("nie udało się sparsować haseł\n")
	}
	req, err := http.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Nie udało się utworzyć żądania HTTP: %v", err))
	}
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var userInstance account.LoginResult
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&userInstance)
	if err != nil {
		return nil, errors.New("nie udało się odczytać response z loginu")
	}

	return &userInstance, nil
}

func testHurt(ean string, token string, router *gin.Engine, t *testing.T, hurt hurtownie.HurtName, shouldWork bool, shouldExist bool) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, "/api/takePrice?ean="+ean, nil)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Nie udało się utworzyć żądania HTTP: %v", err))
	}

	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()

	// Uruchamianie testowego routera Gin z fałszywym żądaniem
	router.ServeHTTP(w, req)
	if shouldWork {
		assert.Equal(t, http.StatusOK, w.Code)
	} else {
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		return false, errors.New("powinno nie dać się przejść do zapytań beaz tokena\n")
	}

	decoder := json.NewDecoder(w.Body)
	var results []takePrices.SearchResult
	err = decoder.Decode(&results)

	filterResult := func(results []takePrices.SearchResult, hurt hurtownie.HurtName) (*takePrices.SearchResult, error) {
		for _, res := range results {
			if res.HurtName == hurt {
				return &res, nil
			}
		}
		return nil, errors.New("there is no matching hurt\n")
	}

	switch hurt {
	case hurtownie.Sot:
		result, err := filterResult(results, hurtownie.Sot)
		assert.Equal(t, err, nil)
		if result.Result == nil {
			return false, errors.New("resIsNil")
		}

		var sotResult hurtownie.SotAndSpecjalResponse

		if resultData, err := json.Marshal(result.Result); err == nil {
			if json.Unmarshal(resultData, &sotResult) != nil {
				return false, errors.New("cant parse response")
			}
		} else {
			return false, err
		}

		assert.NotEqual(t, sotResult.Pozycje[0], nil)
		assert.Contains(t, sotResult.Pozycje[0].Name, "MARS")
	case hurtownie.Specjal: // powinnyśmy dostać żywca
		result, err := filterResult(results, hurtownie.Specjal)
		assert.Equal(t, err, nil)
		if result.Result == nil {
			return false, errors.New("resIsNil")
		}
		var specjalResponse hurtownie.SotAndSpecjalResponse
		if resultData, err := json.Marshal(result.Result); err == nil {
			if json.Unmarshal(resultData, &specjalResponse) != nil {
				return false, errors.New("cant parse response")
			}
		} else {
			return false, err
		}

		assert.NotEqual(t, specjalResponse.Pozycje[0], nil)
		assert.Contains(t, specjalResponse.Pozycje[0].Name, "ŻYWIEC")
	case hurtownie.Tedi: // powinnyśmy dostać żywca
		result, err := filterResult(results, hurtownie.Tedi)
		assert.Equal(t, err, nil)
		if result.Result == nil {
			return false, errors.New("resIsNil")
		}
		var tediResponse tedi.ProductResponse
		if resultData, err := json.Marshal(result.Result); err == nil {
			if json.Unmarshal(resultData, &tediResponse) != nil {
				return false, errors.New("cant parse response")
			}
			if !shouldExist && tediResponse.Count == 0 {
				return false, errors.New("product should not exist")
			}
		} else {
			return false, err
		}
		assert.NotEqual(t, tediResponse.Results[0], nil)
		assert.Contains(t, tediResponse.Results[0].Name, "ŻYWIEC")
	case hurtownie.Eurocash: // powinnyśmy dostać żywca
		result, err := filterResult(results, hurtownie.Eurocash)
		assert.Equal(t, err, nil)
		var eurocashResponse eurocash.EurocashResponse
		if resultData, err := json.Marshal(result.Result); err == nil {
			if json.Unmarshal(resultData, &eurocashResponse) != nil {
				return false, errors.New("cant parse response")
			}
			if !shouldExist && eurocashResponse.Success == true {
				return false, errors.New("product should not exist")
			}
		} else {
			return false, err
		}
		assert.Equal(t, eurocashResponse.Success, true)
		assert.Contains(t, eurocashResponse.Data.Items[0].Name, "ŻYWIEC")
	}

	return true, nil
}
