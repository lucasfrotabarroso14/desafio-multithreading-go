package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Address struct {
	Street     string `json:"street,omitempty"`
	Logradouro string `json:"logradouro,omitempty"`
	City       string `json:"city,omitempty"`
	Localidade string `json:"localidade,omitempty"`
	State      string `json:"state,omitempty"`
	UF         string `json:"uf,omitempty"`
	Bairro     string `json:"neighborhood,omitempty"`
	API        string `json:"api,omitempty"`
}

func GetAdressFromBrasilAPI(cep string, ch chan<- *Address) {
	url := "https://brasilapi.com.br/api/cep/v1/" + cep
	res, err := http.Get(url)
	if err != nil {
		ch <- nil
		return

	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ch <- nil
		return

	}
	var address Address
	err = json.Unmarshal(body, &address)
	if err != nil {
		ch <- nil
		return
	}
	address.API = "BRASIL API"
	ch <- &address

}

func GetAddressFromViaCEP(cep string, ch chan<- *Address) {
	url := "http://viacep.com.br/ws/" + cep + "/json/"
	res, err := http.Get(url)
	if err != nil {
		ch <- nil
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ch <- nil
		return
	}
	var address Address
	err = json.Unmarshal(body, &address)
	if err != nil {
		ch <- nil
		return
	}
	address.API = "VIA CEP"
	ch <- &address
}

func main() {
	cep := "60135280"
	ch := make(chan *Address, 2)
	go GetAdressFromBrasilAPI(cep, ch)
	go GetAddressFromViaCEP(cep, ch)

	var address *Address

	select {
	case address = <-ch:
		if address != nil {
			if address.API == "BRASIL API" {
				fmt.Println("Endereço obtido da BrasilAPI:")
				fmt.Printf("Rua: %s\nCidade: %s\nEstado: %s\nBairro: %s\n", address.Street, address.City, address.State, address.Bairro)

			} else if address.API == "VIA CEP" {
				fmt.Println("Endereço obtido da ViaCEP:")

				fmt.Printf("Rua: %s\nCidade: %s\nEstado: %s\nBairro: %s\n", address.Logradouro, address.Localidade, address.UF, address.Bairro)
			} else {
				fmt.Println("Nenhuma resposta válida recebida.")
			}
		} else {
			fmt.Println("Erro ao obter o endereço.")
		}
	}

}
