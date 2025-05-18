package tedi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"optimaHurt/hurtownie"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ICOM: aby dodać cos do koszyka wystarczy posiadać Ean i ile produktów chcemy dodać --  proste w chuj
type Tedi struct {
	Token hurtownie.SotAndSpecjalTokenResponse
}

func (t *Tedi) CheckToken(client *http.Client) bool {
	return hurtownie.CheckExpDateJwt(t.Token.AccessToken)
}

func (t *Tedi) RefreshToken(client *http.Client) bool {
	return false
} // do naprawy, nie działa

type Creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (t *Tedi) TakeToken(login, password string, client *http.Client) bool {
	url := "https://tedi-ws.ampli-solutions.com/auth/?session_id=" + RandomString(128)
	body := Creds{
		login,
		password,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return false
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return false
	}

	var tediTokenResponse hurtownie.SotAndSpecjalTokenResponse
	defer resp.Body.Close()
	responseReaderJson := json.NewDecoder(resp.Body)
	err = responseReaderJson.Decode(&tediTokenResponse)
	if err != nil {
		return false
	}
	t.Token = tediTokenResponse
	return true
}

func (t *Tedi) SearchProduct(Ean string, client *http.Client) (interface{}, error) {
	req, err := http.NewRequest("GET", "https://tedi-ws.ampli-solutions.com/product-search/?limit=12&search="+Ean+"&offset=0", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "PL")
	req.Header.Add("Amper_app_name", "B2B")
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("child-customer-id", "3620874")
	req.Header.Add("Origin", "https://tedi.kd-24.pl")
	req.Header.Add("Referer", "https://tedi.kd-24.pl")
	req.Header.Add("Authorization", "Bearer "+t.Token.AccessToken)
	req.Header.Add("User-Agent", "cUrl")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("err := %v\n", err)
		return nil, err
	}
	if resp.StatusCode == 403 {
		time.Sleep(time.Second)
		return t.SearchProduct(Ean, client)
	}
	var serverResponse ProductResponse
	defer resp.Body.Close()
	responseReaderJson := json.NewDecoder(resp.Body)
	err = responseReaderJson.Decode(&serverResponse)
	if err != nil {
		fmt.Printf("err %v\n", err)
		return nil, err
	}
	return serverResponse, nil
}

func (t *Tedi) SearchMany(list hurtownie.WishList, client *http.Client) ([]hurtownie.SearchManyProducts, error) {
	ch := make(chan hurtownie.SearchManyProducts)
	var wg sync.WaitGroup
	for _, i := range list.Items {
		wg.Add(1)
		go func(wg *sync.WaitGroup, ch chan<- hurtownie.SearchManyProducts, ean string) {
			defer wg.Done()
			res, err := t.SearchProduct(i.Ean, client)
			if err != nil {
				ch <- hurtownie.SearchManyProducts{
					Item: nil,
					Ean:  ean,
				}
				return
			}
			ch <- hurtownie.SearchManyProducts{
				Item: res,
				Ean:  ean,
			}
		}(&wg, ch, i.Ean)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	result := make([]hurtownie.SearchManyProducts, len(list.Items))
	i := 0
	for x := range ch {
		result[i] = x
		i++
	}
	return result, nil
}

func (t *Tedi) AddToCart(list hurtownie.WishList, client *http.Client) bool {

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	mapEanIntoQuantity := make(map[string]int)
	for _, i := range list.Items {
		mapEanIntoQuantity[i.Ean] = i.Amount
	}

	data := ""
	prodForTedi := make([]hurtownie.Items, 0)
	// wyciąganie danych z wishListy na KCfirme, gdize nie liczy się ani 3 ani 4 kolumna
	for _, i := range list.Items {
		if i.HurtName == hurtownie.Tedi {
			data += fmt.Sprintf("%v;%v;1;a\n", i.Ean, i.Amount)
			prodForTedi = append(prodForTedi, i)
		}
	}
	if len(prodForTedi) == 0 {
		return true
	}

	part, err := writer.CreateFormFile("order_file", "order.txt")

	if err != nil {
		fmt.Println(err)
		return false
	}

	if _, err = io.Copy(part, strings.NewReader(data)); err != nil {
		fmt.Println(err)
		return false
	}

	if writer.WriteField("order_format", "kcfirma") != nil {
		fmt.Println(err)
		return false
	}

	if writer.Close() != nil {
		fmt.Println(err)
		return false
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "https://tedi-ws.ampli-solutions.com/create-order-from-file/", &body)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Set the headers
	setHeader := func(req *http.Request) {
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Accept-Language", "PL")
		req.Header.Add("Amper_app_name", "B2B")
		req.Header.Add("Origin", "https://tedi.kd-24.pl")
		req.Header.Add("Referer", "https://tedi.kd-24.pl")
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
		req.Header.Set("Accept-Language", "pl")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Sec-Fetch-Mode", "no-cors")
		req.Header.Set("Sec-Fetch-Site", "cross-site")
		req.Header.Set("AMPER_APP_NAME", "B2B")
		req.Header.Set("Pragma", "no-cache")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authorization", "Bearer "+t.Token.AccessToken)
	}
	setHeader(req)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()

	var responseWithId struct {
		Data struct {
			CartEntries []struct {
				Id        int `json:"product"`
				BestPrice struct {
					Price float64 `json:"price"`
				} `json:"best_price"`
			} `json:"cart_entries"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&responseWithId); err != nil {
		fmt.Printf("error in tedi get Id := %v", err)
		return false
	}
	items := ""
	for _, i := range responseWithId.Data.CartEntries {
		items += strconv.Itoa(i.Id) + ","
	}
	items = items[:len(items)-1] // usuwamy przecinek z końca

	req, err = http.NewRequest("GET", "https://tedi-ws.ampli-solutions.com/product-list/?ids="+items, nil)
	if err != nil {
		fmt.Println(err)
		return false
	}
	setHeader(req)

	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("%v", err)
		return false
	}

	defer resp.Body.Close()

	var resultWithPrice []struct {
		Id              int     `json:"id"`
		Ean             string  `json:"ean"`
		FinalPrice      float64 `json:"final_price"`
		FinalPriceGross float64 `json:"final_price_gross"`
		BasePrice       float64 `json:"default_price"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&resultWithPrice); err != nil {
		fmt.Printf("error with parsing prices in tedi cart := %v\n", err)
		return false
	}

	type Line struct {
		BasePrice       float64 `json:"base_price_net"`
		FinalPriceGross float64 `json:"final_price_gross"`
		FinalPriceNet   float64 `json:"final_price_net"`
		ProductId       int     `json:"product_id"`
		Quantity        int     `json:"quantity"`
		StockLocation   int     `json:"source_stock_location_id"`
	}

	type Lines struct {
		Line Line `json:"line"`
	}

	type AddToCartPayload struct {
		Lines []Lines `json:"lines"`
	}

	// ICOM: ogarnąć tworzenie tych payloadów itd, chyba powinien być fajrant
	payload := AddToCartPayload{Lines: make([]Lines, 0)}
	for _, i := range resultWithPrice {
		payload.Lines = append(payload.Lines, Lines{
			Line: Line{
				BasePrice:       i.BasePrice,
				FinalPriceNet:   i.FinalPrice,
				FinalPriceGross: i.FinalPriceGross,
				ProductId:       i.Id,
				Quantity:        mapEanIntoQuantity[i.Ean],
				StockLocation:   600016,
			},
		})
	}

	payloadParsed, err := json.Marshal(payload)

	if err != nil {
		fmt.Printf("%v", err)
		return false
	}
	req, err = http.NewRequest("POST", "https://tedi-ws.ampli-solutions.com/cart-lines/?cart_id=31338", bytes.NewBuffer(payloadParsed))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	setHeader(req)
	resp, err = client.Do(req)

	if err != nil {
		fmt.Printf("%v\n", err)
		return false
	}
	if resp.StatusCode != 201 {
		defer resp.Body.Close()
		bdy, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%v\n", err)
			return false
		}
		fmt.Printf("%s\n", bdy)
		return false
	}
	// tutaj jako response dostajemy dane oraz id produktów, na ich podstawie jesteśmy w stanie dostać ceny produktów i przesłać je do koszyka po odpowiedniej cenie

	return true
}

func (t *Tedi) GetName() hurtownie.HurtName {
	return hurtownie.Tedi
}
