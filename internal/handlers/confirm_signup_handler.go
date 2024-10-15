package handlers

import (
	"cognito-lambda-handler/internal/cognito"
	"encoding/json"
	"log"
	"net/http"
)

type ConfirmSignUpRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func ConfirmSignUpHandler(w http.ResponseWriter, r *http.Request, cognitoService *cognito.Service) {
	var req ConfirmSignUpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// ユーザー確認処理
	err = cognitoService.ConfirmSignUp(req.Email, req.Code)
	if err != nil {
		log.Printf("Error confirming sign up for user %s: %v", req.Email, err)
		http.Error(w, "Failed to confirm user", http.StatusInternalServerError)
		return
	}

	// 成功メッセージの返却
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "User confirmed"}); err != nil {
		log.Printf("Error encoding response for user %s: %v", req.Email, err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
