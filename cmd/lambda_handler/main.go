package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cognito-lambda-handler/internal/cognito"
	"github.com/aws/aws-lambda-go/lambda"
)

// SignUpEvent はLambda関数に渡されるイベントの構造体です。
type SignUpEvent struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	GivenName   string `json:"given_name"`
	FamilyName  string `json:"family_name"`
}

func Handler(ctx context.Context, event SignUpEvent) (string, error) {
	// Cognitoクライアントの初期化
	clientID := os.Getenv("AWS_COGNITO_CLIENT_ID")
	if clientID == "" {
		log.Fatal("AWS_COGNITO_CLIENT_ID is not set")
	}

	// Cognitoサービスのインスタンスを作成
	cognitoService, err := cognito.NewCognitoService(clientID)
	if err != nil {
		return "", fmt.Errorf("failed to initialize Cognito service: %v", err)
	}

	// サインアップ処理の実行
	err = cognitoService.SignUp(event.Email, event.Password, event.PhoneNumber, event.GivenName, event.FamilyName)
	if err != nil {
		return "", fmt.Errorf("failed to sign up user: %v", err)
	}

	// 成功時のメッセージを返す
	return fmt.Sprintf("User %s signed up successfully", event.Email), nil
}

func main() {
	// Lambda関数のハンドラを起動
	lambda.Start(Handler)
}
