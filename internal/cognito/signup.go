package cognito

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func (s *Service) SignUp(email, password, phoneNumber, givenName, familyName string) error {
	secretHash, err := generateSecretHash(email, s.clientID)
	if err != nil {
		return fmt.Errorf("failed to generate secret hash: %v", err)
	}

	input := &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(s.clientID),
		SecretHash: aws.String(secretHash),
		Username:   aws.String(email),
		Password:   aws.String(password),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: aws.String(email)},
			{Name: aws.String("phone_number"), Value: aws.String(phoneNumber)},
			{Name: aws.String("given_name"), Value: aws.String(givenName)},
			{Name: aws.String("family_name"), Value: aws.String(familyName)},
		},
	}

	_, err = s.client.SignUp(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to sign up user: %w", err)
	}

	return nil
}
