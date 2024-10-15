package handlers

import (
	"cognito-lambda-handler/internal/cognito"
	"encoding/json"
	"log"
	"net/http"
)

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignInHandler(w http.ResponseWriter, r *http.Request, cognitoService *cognito.Service) {
	var req SignInRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// サインイン処理
	token, err := cognitoService.SignIn(req.Email, req.Password)
	if err != nil {
		log.Printf("Error signing in user %s: %v", req.Email, err)
		http.Error(w, "Failed to sign in user", http.StatusUnauthorized)
		return
	}

	// アクセストークンの返却
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		log.Printf("Error encoding response for user %s: %v", req.Email, err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
