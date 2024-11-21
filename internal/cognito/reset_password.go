package cognito

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

func (s *Service) ResetPassword(email, confirmationCode, newPassword string) error {
	secretHash, err := generateSecretHash(email, s.clientId)
	if err != nil {
		return fmt.Errorf("failed to generate secret hash: %v", err)
	}

	input := &cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         aws.String(s.clientId),
		SecretHash:       aws.String(secretHash),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(confirmationCode),
		Password:         aws.String(newPassword),
	}

	_, err = s.client.ConfirmForgotPassword(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}

	return nil
}
