package cognito

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"os"
)

type Service struct {
	client   *cognitoidentityprovider.Client
	clientID string
}

func NewCognitoService(clientID string) (*Service, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load SDK config: %w", err)
	}

	cognitoClient := cognitoidentityprovider.NewFromConfig(cfg)

	return &Service{
		client:   cognitoClient,
		clientID: clientID,
	}, nil
}

func generateSecretHash(email string, clientID string) (string, error) {
	clientSecret := os.Getenv("AWS_COGNITO_CLIENT_SECRET")
	if clientSecret == "" {
		return "", fmt.Errorf("client secret is not set")
	}

	h := hmac.New(sha256.New, []byte(clientSecret))
	h.Write([]byte(email + clientID))
	secretHash := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return secretHash, nil
}
