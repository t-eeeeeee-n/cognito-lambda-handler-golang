package handlers

import (
	"cognito-lambda-handler/internal/cognito"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/aws/smithy-go"
)

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignInHandler(w http.ResponseWriter, r *http.Request, cognitoService *cognito.Service) {
	if cognitoService == nil {
		http.Error(w, "Cognito service is not initialized", http.StatusInternalServerError)
		return
	}

	var req SignInRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := cognitoService.SignIn(req.Email, req.Password)
	//authResult, err := cognitoService.SignIn(req.Email, req.Password)
	if err != nil {
		var awsErr smithy.APIError
		if ok := errors.As(err, &awsErr); ok {
			switch awsErr.ErrorCode() {
			case "NotAuthorizedException":
				http.Error(w, "Incorrect username or password", http.StatusUnauthorized)
			case "UserNotFoundException":
				http.Error(w, "User does not exist", http.StatusNotFound)
			case "InvalidParameterException":
				http.Error(w, "Invalid input parameters", http.StatusBadRequest)
			default:
				http.Error(w, "Failed to sign in user", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Failed to sign in user", http.StatusInternalServerError)
		}
		log.Printf("Error signing in user %s: %v", req.Email, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		log.Printf("Error encoding response for user %s: %v", req.Email, err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

	//response := map[string]interface{}{
	//	"accessToken":  *authResult.AccessToken,
	//	"idToken":      *authResult.IdToken,
	//	"refreshToken": authResult.RefreshToken,
	//	"tokenType":    authResult.TokenType,
	//	"expiresIn":    authResult.ExpiresIn,
	//}
	//
	//if err := json.NewEncoder(w).Encode(response); err != nil {
	//	log.Printf("Error encoding response for user %s: %v", req.Email, err)
	//	http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	//}
}
