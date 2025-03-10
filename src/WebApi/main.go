package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Data models (matching the .NET DTOs)

type SaldoDto struct {
	Total      int       `json:"total"`
	Limite     int       `json:"limite"`
	DataExtrato time.Time `json:"data_extrato"`
}

type TransacaoDto struct {
	Valor     int    `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}

type ExtratoDto struct {
	Saldo             SaldoDto       `json:"saldo"`
	UltimasTransacoes []TransacaoDto `json:"ultimas_transacoes"`
}

type ClienteDto struct {
	Id     int `json:"id"`
	Limite int `json:"limite"`
	Saldo  int `json:"saldo"`
}

// Global clientes mapping as in the .NET code
var clientes = map[int]int{
	1: 100000,
	2: 80000,
	3: 1000000,
	4: 10000000,
	5: 500000,
}

var dbPool *pgxpool.Pool

func main() {
	// Read connection string from environment variable.
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL env var is not set")
	}

	// Create a pgx pool
	var err error
	dbPool, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Create router and endpoints.
	r := chi.NewRouter()

	r.Get("/healthz", healthzHandler)
	r.Get("/clientes/{id}/extrato", getExtratoHandler)
	r.Post("/clientes/{id}/transacoes", postTransacaoHandler)

	srv := &http.Server{
		Addr:         ":9999",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Starting server on port 9999...")

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// healthzHandler is a basic health check
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// getExtratoHandler implements GET /clientes/{id}/extrato
func getExtratoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	clientId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	limite, exists := clientes[clientId]
	if !exists {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	row := dbPool.QueryRow(ctx, "SELECT * FROM GetSaldoClienteById($1)", clientId)

	var total int
	var dbLimite int
	var dataExtrato time.Time
	var transacoesJSON []byte

	err = row.Scan(&total, &dbLimite, &dataExtrato, &transacoesJSON)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Extrato not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Error reading from database: %v", err), http.StatusInternalServerError)
		return
	}

	// Unmarshal the JSON array of transactions.
	var transacoes []TransacaoDto
	if err := json.Unmarshal(transacoesJSON, &transacoes); err != nil {
		// If unmarshal fails, log and use empty slice.
		log.Printf("Failed to unmarshal transacoes JSON: %v", err)
		transacoes = []TransacaoDto{}
	}

	// Return the result.
	extrato := ExtratoDto{
		Saldo: SaldoDto{
			Total:      total,
			Limite:     dbLimite,
			DataExtrato: dataExtrato,
		},
		UltimasTransacoes: transacoes,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(extrato); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding result: %v", err), http.StatusInternalServerError)
	}
}

// postTransacaoHandler implements POST /clientes/{id}/transacoes
func postTransacaoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	clientId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	limite, exists := clientes[clientId]
	if !exists {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	// Parse incoming JSON
	var transacao TransacaoDto
	if err := json.NewDecoder(r.Body).Decode(&transacao); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate transaction input
	if !isTransacaoValid(transacao) {
		http.Error(w, "Invalid transaction data", http.StatusUnprocessableEntity)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	row := dbPool.QueryRow(ctx, "SELECT InsertTransacao($1, $2, $3, $4)", clientId, transacao.Valor, transacao.Tipo, transacao.Descricao)
	var updatedSaldo int
	err = row.Scan(&updatedSaldo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error inserting transaction: %v", err), http.StatusInternalServerError)
		return
	}

	clienteDto := ClienteDto{
		Id:     clientId,
		Limite: limite,
		Saldo:  updatedSaldo,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(clienteDto); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding result: %v", err), http.StatusInternalServerError)
	}
}

// isTransacaoValid validates the transaction just like the .NET code.
func isTransacaoValid(transacao TransacaoDto) bool {
	tipoC := "c"
	tipoD := "d"

	if transacao.Tipo != tipoC && transacao.Tipo != tipoD {
		return false
	}
	if transacao.Descricao == "" {
		return false
	}
	if len(transacao.Descricao) > 10 {
		return false
	}
	if transacao.Valor <= 0 {
		return false
	}
	return true
}