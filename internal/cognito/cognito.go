package cognito

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

// Service はCognitoの操作を提供する構造体です。
type Service struct {
	client   *cognitoidentityprovider.Client
	clientID string
}

// NewCognitoService は新しいCognitoサービスを初期化します。
func NewCognitoService(clientID string) (*Service, error) {
	// AWS SDKの設定を読み込む
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load SDK config: %w", err)
	}

	// Cognitoクライアントの作成
	cognitoClient := cognitoidentityprovider.NewFromConfig(cfg)

	return &Service{
		client:   cognitoClient,
		clientID: clientID,
	}, nil
}

func generateSecretHash(email string, clientID string) (string, error) {
	clientSecret := os.Getenv("AWS_COGNITO_CLIENT_SECRET") // クライアントシークレットを環境変数から取得
	if clientSecret == "" {
		return "", fmt.Errorf("client secret is not set")
	}

	h := hmac.New(sha256.New, []byte(clientSecret))
	h.Write([]byte(email + clientID))
	secretHash := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return secretHash, nil
}

// SignUp ユーザーをCognitoにサインアップさせます。
func (s *Service) SignUp(email, password, phoneNumber, givenName, familyName string) error {
	secretHash, err := generateSecretHash(email, s.clientID)
	if err != nil {
		return fmt.Errorf("failed to generate secret hash: %v", err)
	}

	input := &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(s.clientID),
		SecretHash: aws.String(secretHash),
		Username:   aws.String(email), // EメールをUsernameとして使う
		Password:   aws.String(password),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
			{
				Name:  aws.String("phone_number"),
				Value: aws.String(phoneNumber), // 電話番号
			},
			{
				Name:  aws.String("given_name"),
				Value: aws.String(givenName), // 名
			},
			{
				Name:  aws.String("family_name"),
				Value: aws.String(familyName), // 姓
			},
		},
	}

	_, err = s.client.SignUp(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to sign up user: %w", err)
	}

	return nil
}
