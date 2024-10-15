package handlers

import (
	"cognito-lambda-handler/internal/cognito"
	"encoding/json"
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
	if cognitoService == nil {
		http.Error(w, "Cognito service is not initialized", http.StatusInternalServerError)
		return
	}

	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := cognitoService.SignUp(req.Email, req.Password, req.PhoneNumber, req.GivenName, req.FamilyName); err != nil {
		http.Error(w, "Failed to sign up user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Sign up successful"})
}
