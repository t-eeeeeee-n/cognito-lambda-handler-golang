package cognito

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

func (s *Service) ForgotPassword(email string) error {
	secretHash, err := generateSecretHash(email, s.clientId)
	if err != nil {
		return fmt.Errorf("failed to generate secret hash: %v", err)
	}

	input := &cognitoidentityprovider.ForgotPasswordInput{
		ClientId:   aws.String(s.clientId),
		SecretHash: aws.String(secretHash),
		Username:   aws.String(email),
	}

	_, err = s.client.ForgotPassword(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to request password reset: %w", err)
	}

	return nil
}
