package hurtownie

import (
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"strings"
	"time"
)

func GenerateUUID() string {
	var D [36]byte
	A := "0123456789abcdef"

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	for B := 0; B < 36; B++ {
		D[B] = A[rand.Intn(16)]
	}
	D[14] = '4'
	D[19] = A[(D[19]&0x3)|0x8]
	D[8] = '-'
	D[13] = '-'
	D[18] = '-'
	D[23] = '-'
	C := strings.Join([]string{string(D[:])}, "")
	return C
}

func CheckExpDateJwt(hurtToken string) bool {
	token, _, err := jwt.NewParser().ParseUnverified(hurtToken, jwt.MapClaims{})
	if err != nil {
		return false
	}

	// Wyciąganie informacji o dacie wygaśnięcia (exp)
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			// Konwertowanie czasu wygaśnięcia z Unix Timestamp na time.Time
			expirationTime := time.Unix(int64(exp)-1 /*odjęcie sekundy ważności tokenu aby mieć margines błędu przetwarzania requesta*/, 0)
			// Sprawdzenie, czy token jest już przeterminowany
			if !time.Now().After(expirationTime) {
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	} else {
		return false
	}
}

type SotAndSpecjalTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SotAndSpecjalResponse struct {
	DataOferty time.Time `json:"dataOferty"`
	Pozycje    []struct {
		Name                string  `json:"nazwa"`
		IlOpkZb             int     `json:"ilOpkZb"`
		CenaNettoOstateczna float64 `json:"cenaNettoOstateczna"`
	} `json:"pozycje"`
	CountPozycji int `json:"countPozycji"`
}
