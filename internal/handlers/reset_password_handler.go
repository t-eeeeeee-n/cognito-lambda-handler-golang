package handlers

import (
	"cognito-lambda-handler/internal/cognito"
	"encoding/json"
	"log"
	"net/http"
)

type ResetPasswordRequest struct {
	Email       string `json:"email"`
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request, cognitoService *cognito.Service) {
	var req ResetPasswordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// パスワードリセットの確認
	err = cognitoService.ResetPassword(req.Email, req.Code, req.NewPassword)
	if err != nil {
		log.Printf("Error resetting password for user %s: %v", req.Email, err)
		http.Error(w, "Failed to reset password", http.StatusInternalServerError)
		return
	}

	// 成功メッセージの返却
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Password reset successful"}); err != nil {
		log.Printf("Error encoding response for user %s: %v", req.Email, err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
