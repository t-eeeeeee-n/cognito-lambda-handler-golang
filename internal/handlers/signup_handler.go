package handlers

import (
	"cognito-lambda-handler/internal/cognito"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/aws/smithy-go"
)

type SignUpRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	GivenName   string `json:"given_name"`
	FamilyName  string `json:"family_name"`
}

func SignUpHandler(w http.ResponseWriter, r *http.Request, cognitoService *cognito.Service) {
	// Cognitoサービスの初期化確認
	if cognitoService == nil {
		http.Error(w, "Cognito service is not initialized", http.StatusInternalServerError)
		return
	}

	var req SignUpRequest

	// リクエストのデコード
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// サインアップ処理の呼び出し
	err = cognitoService.SignUp(req.Email, req.Password, req.PhoneNumber, req.GivenName, req.FamilyName)
	if err != nil {
		var awsErr smithy.APIError
		// AWSエラーが発生した場合の処理
		if ok := errors.As(err, &awsErr); ok {
			switch awsErr.ErrorCode() {
			case "UsernameExistsException":
				// メールアドレスの重複エラー
				http.Error(w, "User with this email already exists", http.StatusConflict)
			case "InvalidParameterException":
				// 不正なパラメータが渡された場合
				http.Error(w, "Invalid input parameters", http.StatusBadRequest)
			case "CodeMismatchException":
				// 認証コードが間違っている場合
				http.Error(w, "Incorrect verification code", http.StatusBadRequest)
			case "LimitExceededException":
				// リクエストのレート制限を超えた場合
				http.Error(w, "Request limit exceeded", http.StatusTooManyRequests)
			case "NotAuthorizedException":
				// ユーザーが認証されていない場合
				http.Error(w, "Not authorized", http.StatusUnauthorized)
			case "ResourceNotFoundException":
				// リソースが見つからない場合
				http.Error(w, "Resource not found", http.StatusNotFound)
			default:
				// その他のエラー
				http.Error(w, "Failed to sign up user", http.StatusInternalServerError)
			}
		} else {
			// その他のエラー
			http.Error(w, "Failed to sign up user", http.StatusInternalServerError)
		}
		log.Printf("Error signing up user %s: %v", req.Email, err)
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
