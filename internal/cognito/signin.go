package cognito

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"time"
)

// SignIn 全てreturnする場合は*types.AuthenticationResultType
func (s *Service) SignIn(email, password string) (string, error) {
	// SRPオブジェクトの作成
	srp, err := NewCognitoSRP(email, password, s.poolId, s.clientId, s.clientSecret)
	if err != nil {
		return "", fmt.Errorf("failed to create SRP object: %w", err)
	}

	// InitiateAuthリクエストを作成
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow:       types.AuthFlowTypeUserSrpAuth,
		AuthParameters: srp.GetAuthParams(),
		ClientId:       aws.String(s.clientId),
	}

	// InitiateAuthの呼び出し
	output, err := s.client.InitiateAuth(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to initiate auth: %w", err)
	}

	// チャレンジレスポンスを処理
	challengeResponse, err := srp.PasswordVerifierChallenge(output.ChallengeParameters, time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to calculate challenge response: %w", err)
	}

	// チャレンジに応答
	respondInput := &cognitoidentityprovider.RespondToAuthChallengeInput{
		ChallengeName:      output.ChallengeName,
		ChallengeResponses: challengeResponse,
		ClientId:           aws.String(s.clientId),
	}

	authResult, err := s.client.RespondToAuthChallenge(context.TODO(), respondInput)
	if err != nil {
		return "", fmt.Errorf("failed to respond to auth challenge: %w", err)
	}

	return *authResult.AuthenticationResult.AccessToken, nil
}
