package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Endereco struct {
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

// Thread 1
func main() {
	ch1 := make(chan Endereco)
	ch2 := make(chan Endereco)
	defer close(ch1)
	defer close(ch2)

	cep := "29946590"
	go consultarCepCdn(cep, ch1)
	go consultarCepSite(cep, ch2)

	select {
	case msg := <-ch1:
		fmt.Printf("consultarCepCdn: %v\n", msg)
	case msg := <-ch2:
		fmt.Printf("consultarCepSite: %v\n", msg)
	case <-time.After(1 * time.Second):
		fmt.Println("timeout")
	}
}

func consultarCepSite(cep string, ch chan<- Endereco) {
	endereco, err := obterEndereco("http://viacep.com.br/ws/" + cep + "/json")
	if err != nil {
		fmt.Printf("erro ao obter o cep via cdn: %v", err)
		return
	}
	ch <- endereco
}

func consultarCepCdn(cep string, ch chan<- Endereco) {
	endereco, err := obterEndereco("https://cdn.apicep.com/file/apicep/" + cep + ".json")
	if err != nil {
		fmt.Printf("erro ao obter o cep via cdn: %v", err)
		return
	}
	ch <- endereco
}

func obterEndereco(url string) (Endereco, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Endereco{}, err
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return Endereco{}, err
	}
	var endereco Endereco
	err = json.Unmarshal(content, &endereco)
	if err != nil {
		return Endereco{}, err
	}
	return endereco, nil
}
