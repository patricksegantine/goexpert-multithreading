package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Endereco struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
}

// Thread 1
func main() {
	ch1 := make(chan Endereco)
	ch2 := make(chan Endereco)
	defer close(ch1)
	defer close(ch2)

	cep := "29946-590"
	go consultarCepSite(strings.Replace(cep, "-", "", -1), ch2)
	go consultarCepCdn(cep, ch1)

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
	content, err := obterEndereco(fmt.Sprintf("http://viacep.com.br/ws/%s/json", cep))
	if err != nil {
		fmt.Printf("erro ao obter o cep via cdn: %v", err)
		return
	}
	var endereco Endereco
	err = json.Unmarshal(content, &endereco)
	if err != nil {
		return
	}
	ch <- endereco
}

func consultarCepCdn(cep string, ch chan<- Endereco) {
	content, err := obterEndereco(fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", cep))
	if err != nil {
		fmt.Printf("erro ao obter o cep via cdn: %v", err)
		return
	}

	var tempEnd map[string]interface{}
	err = json.Unmarshal(content, &tempEnd)
	if err != nil {
		return
	}
	endereco := Endereco{
		Cep:        tempEnd["code"].(string),
		Logradouro: tempEnd["address"].(string),
		Bairro:     tempEnd["district"].(string),
		Localidade: tempEnd["city"].(string),
		Uf:         tempEnd["state"].(string),
	}
	ch <- endereco
}

func obterEndereco(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}
