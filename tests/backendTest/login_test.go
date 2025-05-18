package backendTest

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"optimaHurt/constAndVars"
	"optimaHurt/endpoints"
	"optimaHurt/endpoints/account"
	"optimaHurt/user"
	"os"
	"testing"
)

func TestMaker(t *testing.T) {

	if constAndVars.DbConnect != nil {
		return
	}

	connection := endpoints.ConnectToDB(os.Getenv("CONNECTION_STRING"))
	defer connection.Disconnect(context.Background())

	r := gin.Default()
	r.POST("/api/login", account.Login)
	r.GET("/api/checkCookie", account.TestCookie)

	LoginSimple(t, r)
	LoginChain(t, r)

}

func LoginSimple(t *testing.T, r *gin.Engine) {
	w := httptest.NewRecorder()
	body := user.LoginBodyData{Username: "kholin", Password: "nice"}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(data))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("błąd na poprawne co kurwa")
		return
	}
}

func LoginChain(t *testing.T, r *gin.Engine) {

	w := httptest.NewRecorder()
	body := user.LoginBodyData{Username: "kholin", Password: "nice"}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(data))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("błąd na poprawne dane")
		return
	}
	var response struct {
		ShieldedToken string `json:"token"`
	}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("error przy odczycie body := %v", err)
		return
	}
	req, _ = http.NewRequest("GET", "/api/checkCookie", nil)
	req.Header.Set("Authorization", response.ShieldedToken)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("zwrucono nieprawidłowy token, checkCookei się nie powiodło")
		return
	}

}
