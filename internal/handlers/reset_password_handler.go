package handlers

import (
	"cognito-lambda-handler/internal/cognito"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/aws/smithy-go"
)

type ResetPasswordRequest struct {
	Email       string `json:"email"`
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request, cognitoService *cognito.Service) {
	if cognitoService == nil {
		http.Error(w, "Cognito service is not initialized", http.StatusInternalServerError)
		return
	}

	var req ResetPasswordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = cognitoService.ResetPassword(req.Email, req.Code, req.NewPassword)
	if err != nil {
		var awsErr smithy.APIError
		if ok := errors.As(err, &awsErr); ok {
			switch awsErr.ErrorCode() {
			case "CodeMismatchException":
				http.Error(w, "Invalid verification code", http.StatusBadRequest)
			case "ExpiredCodeException":
				http.Error(w, "Verification code has expired", http.StatusBadRequest)
			default:
				http.Error(w, "Failed to reset password", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Failed to reset password", http.StatusInternalServerError)
		}
		log.Printf("Error resetting password for user %s: %v", req.Email, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Password reset successful"}); err != nil {
		log.Printf("Error encoding response for user %s: %v", req.Email, err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
