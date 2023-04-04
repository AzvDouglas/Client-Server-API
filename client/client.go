package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	// Definindo o timeout máximo de 300ms para receber o resultado do server.go
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Fazendo a requisição para o endpoint /cotacao no server.go
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		fmt.Println("Erro ao criar requisição:", err)
		return
	}

	// Fazendo a requisição e tratando possíveis erros
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Erro ao fazer requisição:", err)
		return
	}
	defer resp.Body.Close()

	// Lendo a resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler resposta:", err)
		return
	}

	// Verificando se o status code é 200
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Status code inválido:", resp.StatusCode)
		fmt.Println("Resposta:", string(body))
		return
	}
	//ping := string(body)
	//fmt.Println(ping)

	// Parseando o JSON da resposta
	var cotacao struct {
		Bid string `json:"bid"`
	}
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		fmt.Println("Erro ao parsear resposta:", err)
		return
	}

	// Escrevendo o resultado no arquivo cotacao.txt
	err = os.WriteFile("cotacao.txt", []byte(fmt.Sprintf("Dólar:\n$1 = R$%s", cotacao.Bid)), 0644)
	if err != nil {
		fmt.Println("Erro ao escrever no arquivo:", err)
		return
	}

	fmt.Println("Cotação salva com sucesso!")
}
