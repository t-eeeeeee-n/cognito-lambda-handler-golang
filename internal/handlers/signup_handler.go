package handlers

import (
	"cognito-lambda-handler/internal/cognito"
	"encoding/json"
	"log"
	"net/http"
)

type SignUpRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	GivenName   string `json:"given_name"`
	FamilyName  string `json:"family_name"`
}

func SignUpHandler(w http.ResponseWriter, r *http.Request, cognitoService *cognito.Service) {
	var req SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// サインアップ処理の呼び出し
	err = cognitoService.SignUp(req.Email, req.Password, req.PhoneNumber, req.GivenName, req.FamilyName)
	if err != nil {
		log.Printf("Error signing up user %s: %v", req.Email, err)
		http.Error(w, "Failed to sign up user", http.StatusInternalServerError)
		return
	}

	// 成功メッセージの返却
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Sign up successful"}); err != nil {
		log.Printf("Error encoding response for user %s: %v", req.Email, err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
