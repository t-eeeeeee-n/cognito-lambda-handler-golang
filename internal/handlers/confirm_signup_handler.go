package handlers

import (
	"cognito-lambda-handler/internal/cognito"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/aws/smithy-go"
)

type ConfirmSignUpRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func ConfirmSignUpHandler(w http.ResponseWriter, r *http.Request, cognitoService *cognito.Service) {
	if cognitoService == nil {
		http.Error(w, "Cognito service is not initialized", http.StatusInternalServerError)
		return
	}

	var req ConfirmSignUpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = cognitoService.ConfirmSignUp(req.Email, req.Code)
	if err != nil {
		var awsErr smithy.APIError
		if ok := errors.As(err, &awsErr); ok {
			switch awsErr.ErrorCode() {
			case "CodeMismatchException":
				http.Error(w, "Invalid verification code", http.StatusBadRequest)
			case "ExpiredCodeException":
				http.Error(w, "Verification code expired", http.StatusBadRequest)
			default:
				http.Error(w, "Failed to confirm sign up", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Failed to confirm sign up", http.StatusInternalServerError)
		}
		log.Printf("Error confirming sign up for user %s: %v", req.Email, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "User confirmed"}); err != nil {
		log.Printf("Error encoding response for user %s: %v", req.Email, err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
