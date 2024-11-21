package cognito

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

func (s *Service) ConfirmSignUp(email, confirmationCode string) error {
	secretHash, err := generateSecretHash(email, s.clientId)
	if err != nil {
		return fmt.Errorf("failed to generate secret hash: %v", err)
	}

	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(s.clientId),
		SecretHash:       aws.String(secretHash),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(confirmationCode),
	}

	_, err = s.client.ConfirmSignUp(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to confirm sign up: %w", err)
	}

	return nil
}
