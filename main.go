package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BrasilAPIResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCEPResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	cep := "01153000"
	urlBrasilAPI := "https://brasilapi.com.br/api/cep/v1/" + cep
	urlViaCEP := "http://viacep.com.br/ws/" + cep + "/json/"

	brasilAPIChan := make(chan BrasilAPIResponse)
	viaCEPChan := make(chan ViaCEPResponse)
	errChan := make(chan error, 2)

	go FetchAdressBrasilAPI(urlBrasilAPI, brasilAPIChan, errChan)
	go FetchAdressViaCEP(urlViaCEP, viaCEPChan, errChan)

	select {
	case result := <-brasilAPIChan:
		fmt.Println("A API mais rápida foi BrasilAPI")
		fmt.Printf("%+v\n", result)
	case result := <-viaCEPChan:
		fmt.Println("A API mais rápida foi ViaCEP")
		fmt.Printf("%+v\n", result)
	case <-time.After(time.Second * 1):
		println("timeout")
	}
}

func FetchAdressBrasilAPI(url string, resultChan chan<- BrasilAPIResponse, errChan chan<- error) {
	resp, err := http.Get(url)
	if err != nil {
		errChan <- err
		return
	}

	defer resp.Body.Close()

	var response BrasilAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		errChan <- err
		return
	}

	resultChan <- response
}

func FetchAdressViaCEP(url string, resultChan chan<- ViaCEPResponse, errChan chan<- error) {
	resp, err := http.Get(url)
	if err != nil {
		errChan <- err
		return
	}

	defer resp.Body.Close()

	var response ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		errChan <- err
		return
	}

	resultChan <- response
}
