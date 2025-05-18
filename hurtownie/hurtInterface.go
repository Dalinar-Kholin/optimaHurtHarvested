package hurtownie

import "net/http"

type HurtName int

const (
	Eurocash HurtName = 1 << iota
	Specjal
	Sot
	Tedi
)

type Token string

type IHurt interface {
	CheckToken(client *http.Client) bool
	RefreshToken(client *http.Client) bool
	TakeToken(login, password string, client *http.Client) bool // jeśli udało się pobrać token zwraca true wpp zwraca false
	// te 3  powinny być oddelegowane do obiektu Tokena i robione poprzez delegacje

	SearchProduct(Ean string, client *http.Client) (interface{}, error)
	SearchMany(list WishList, client *http.Client) ([]SearchManyProducts, error)

	AddToCart(list WishList, client *http.Client) bool // tutaj będą potrzebne konkretne instancje itemów z konkretnych hurotwni, czy można dać typy generyczne - wyjebane, wszystkow zwracam na frontend i fajrant

	GetName() HurtName
}

type SearchManyProducts struct {
	Item interface{}
	Ean  string
}

type WishList struct {
	Items []Items `json:"Items"`
}
type Items struct {
	Ean      string   `json:"Ean"`
	Amount   int      `json:"Amount"`
	HurtName HurtName `json:"HurtName,omitempty"`
}

/*
refaktoryzacja na poziomie, zmiany interfu pod posiadnie osobnych obiektów od

 -- zarządzania hasłem
 -- dodawania do koszyka
 -- szukania produktu

wszystko robione poprzez kompozycję

inna sprawa, to rozbicie wszystkiego na mikroserwisy

*/
