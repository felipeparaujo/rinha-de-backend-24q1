package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/felipeparaujo/rinha-de-backend-24q1/api"
)

const url = "http://localhost:9999"
const contentType = "application/json"

func main() {
	log.Println("Transacao")
	transacao()

	log.Println("Extrato")
	extrato()
}

func transacao() {
	marshaled, err := json.Marshal(api.TransacoesRequest{Valor: 1000000000, Tipo: "d", Descricao: "desc"})
	r, err := http.DefaultClient.Post(fmt.Sprintf("%s/clientes/%d/transacoes", url, 1), contentType, bytes.NewBuffer(marshaled))
	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()
	a, err := io.ReadAll(r.Body)
	log.Print(string(a))
	log.Print(r.Status)
	log.Print(err)
}

func extrato() {
	r, err := http.DefaultClient.Get(fmt.Sprintf("%s/clientes/%d/extrato", url, 1))
	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()
	a, err := io.ReadAll(r.Body)
	log.Print(string(a))
	log.Print(r.Status)
	log.Print(err)
}
