package handlers

import (
	"cognito-lambda-handler/internal/cognito"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/aws/smithy-go"
)

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request, cognitoService *cognito.Service) {
	if cognitoService == nil {
		http.Error(w, "Cognito service is not initialized", http.StatusInternalServerError)
		return
	}

	var req ForgotPasswordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = cognitoService.ForgotPassword(req.Email)
	if err != nil {
		var awsErr smithy.APIError
		if ok := errors.As(err, &awsErr); ok {
			switch awsErr.ErrorCode() {
			case "UserNotFoundException":
				http.Error(w, "User not found", http.StatusNotFound)
			case "LimitExceededException":
				http.Error(w, "Password reset limit exceeded", http.StatusTooManyRequests)
			default:
				http.Error(w, "Failed to request password reset", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Failed to request password reset", http.StatusInternalServerError)
		}
		log.Printf("Error requesting password reset for user %s: %v", req.Email, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Password reset requested"}); err != nil {
		log.Printf("Error encoding response for user %s: %v", req.Email, err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
