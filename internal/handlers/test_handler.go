package handlers

import (
	"cognito-lambda-handler/internal/services" // services をインポート
	"encoding/json"
	"log"
	"net/http"
)

// TestHandler - TestService を呼び出し、Hello, World! を返す
func TestHandler(w http.ResponseWriter, r *http.Request) {
	message := services.TestService() // サービス層の関数を呼び出し

	response := map[string]string{"message": message}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
