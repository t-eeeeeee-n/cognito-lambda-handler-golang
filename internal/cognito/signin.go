package cognito

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func (s *Service) SignIn(email, password string) (string, error) {
	secretHash, err := generateSecretHash(email, s.clientID)
	if err != nil {
		return "", fmt.Errorf("failed to generate secret hash: %v", err)
	}

	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME":    email,
			"PASSWORD":    password,
			"SECRET_HASH": secretHash,
		},
		ClientId: aws.String(s.clientID),
	}

	output, err := s.client.InitiateAuth(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to sign in: %w", err)
	}

	return *output.AuthenticationResult.AccessToken, nil
}
