package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	timeoutAPICall = 200 * time.Millisecond
	timeoutDBSave  = 10 * time.Millisecond
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	// Conecta ao banco de dados SQLite
	db, err := sql.Open("sqlite3", "cotacao.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Cria a tabela se ela ainda não existir
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS cotacao (id INTEGER PRIMARY KEY, bid STRING, created_at DATETIME)"); err != nil {
		log.Fatal(err)
	}

	// Cria um servidor HTTP e registra um endpoint para a cotação
	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		// Cria um contexto com timeout para a chamada da API de cotação
		ctx, cancel := context.WithTimeout(r.Context(), timeoutAPICall)
		defer cancel()

		// Realiza a chamada à API de cotação
		resp, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Decodifica o JSON de resposta da API
		var data map[string]Cotacao
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Recupera a cotação do dólar do JSON
		bid := data["USDBRL"].Bid

		// Cria um contexto com timeout para salvar a cotação no banco de dados
		ctx2, cancel2 := context.WithTimeout(ctx, timeoutDBSave)
		defer cancel2()

		// Salva a cotação no banco de dados
		if _, err := db.ExecContext(ctx2, "INSERT INTO cotacao (bid, created_at) VALUES (?, ?)", bid, time.Now()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Monta o JSON de resposta para o cliente
		response := map[string]string{"bid": bid}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			return
		}
	})

	// Inicia o servidor HTTP na porta 8080
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
