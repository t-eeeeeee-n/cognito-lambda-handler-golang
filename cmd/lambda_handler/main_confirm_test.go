package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 確認サインアップのテスト
func TestConfirmSignUpHandler_Success(t *testing.T) {
	requestBody, _ := json.Marshal(map[string]string{
		"email": "tensho.arai@ities-inc.co.jp",
		"code":  "943582", // 実際の確認コードをここに
	})

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/confirm",
		Body:       string(requestBody),
	}

	resp, err := Handler(context.Background(), req)

	if err != nil {
		t.Fatalf("Error calling Lambda handler: %v", err)
	}

	assert.Equal(t, 200, resp.StatusCode, "Expected status code to be 200")
	assert.Contains(t, resp.Body, "User confirmed", "Expected confirmation message")
}

// パスワードリセット確認のテスト
//func TestResetPasswordHandler_Success(t *testing.T) {
//	requestBody, _ := json.Marshal(map[string]string{
//		"email":        "testuser@example.com",
//		"code":         "943582", // 実際の確認コードをここに
//		"new_password": "NewPassword123!",
//	})
//
//	req := events.APIGatewayProxyRequest{
//		HTTPMethod: "POST",
//		Path:       "/reset-password",
//		Body:       string(requestBody),
//	}
//
//	resp, err := Handler(context.Background(), req)
//
//	if err != nil {
//		t.Fatalf("Error calling Lambda handler: %v", err)
//	}
//
//	assert.Equal(t, 200, resp.StatusCode, "Expected status code to be 200")
//	assert.Contains(t, resp.Body, "Password reset successful", "Expected reset successful message")
//}
