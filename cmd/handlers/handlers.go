package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zorjak/tx-parser/parser"
)

type CurrentBlockResponse struct {
	CurrentBlock int `json:"current_block"`
}

type SubscribeBody struct {
	Address string `json:"address"`
}

func CurrentBlockHandler(parser parser.Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("GET /api/current-block")
		message := CurrentBlockResponse{CurrentBlock: 0}
		message.CurrentBlock = parser.GetCurrentBlock()
		respondWithJSON(w, http.StatusOK, message)
	}
}

func SubscribeHandler(parser parser.Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("POST /api/subscribe")
		var body SubscribeBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
		}

		parser.Subscribe(body.Address)
		respondWithJSON(w, http.StatusOK, body)
	}
}

func GetTransactionsHandler(parser parser.Parser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("GET /api/transactions")
		vars := mux.Vars(r)
		address := vars["address"]

		transactions := parser.GetTransactions(address)
		respondWithJSON(w, http.StatusOK, transactions)
	}
}

func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
