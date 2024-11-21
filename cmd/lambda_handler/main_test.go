package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

// ユニークなメールアドレスを生成
func generateUniqueEmail() string {
	return fmt.Sprintf("testuser_%d@example.com", time.Now().UnixNano())
}

// サインアップが成功することを確認
func TestSignUpHandler_Success(t *testing.T) {
	email := generateUniqueEmail()
	requestBody, _ := json.Marshal(map[string]string{
		"email":        email,
		"password":     "Password123!", // パスワードをポリシーに適合させる
		"phone_number": "+1234567890",
		"given_name":   "Test",
		"family_name":  "User",
	})

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/signup",
		Body:       string(requestBody),
	}

	resp, err := Handler(context.Background(), req)

	if err != nil {
		t.Fatalf("Error calling Lambda handler: %v", err)
	}

	assert.Equal(t, 200, resp.StatusCode, "Expected status code to be 200")
	assert.Contains(t, resp.Body, "Sign up successful", "Expected successful signup message")
}

// サインアップ失敗（重複メール）のテスト
func TestSignUpHandler_UsernameExists(t *testing.T) {
	requestBody, _ := json.Marshal(map[string]string{
		"email":        "tensho.arai@ities-inc.co.jp",
		"password":     "Password123!", // 同様に、パスワードを適合させる
		"phone_number": "+1234567890",
		"given_name":   "Existing",
		"family_name":  "User",
	})

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/signup",
		Body:       string(requestBody),
	}

	resp, err := Handler(context.Background(), req)

	if err != nil {
		t.Fatalf("Error calling Lambda handler: %v", err)
	}

	assert.Equal(t, 409, resp.StatusCode, "Expected status code to be 409")
	assert.Contains(t, resp.Body, "User with this email already exists", "Expected email exists error message")
}

// サインインが成功することを確認
func TestSignInHandler_Success(t *testing.T) {
	requestBody, _ := json.Marshal(map[string]string{
		"email":    "testuser@example.com",
		"password": "Password123!",
	})

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/signin",
		Body:       string(requestBody),
	}

	resp, err := Handler(context.Background(), req)

	if err != nil {
		t.Fatalf("Error calling Lambda handler: %v", err)
	}

	assert.Equal(t, 200, resp.StatusCode, "Expected status code to be 200")
	assert.Contains(t, resp.Body, "token", "Expected response to contain token")
}

// サインイン失敗（不正なパスワード）のテスト
func TestSignInHandler_InvalidPassword(t *testing.T) {
	requestBody, _ := json.Marshal(map[string]string{
		"email":    "testuser@example.com",
		"password": "WrongPassword!", // 誤ったパスワード
	})

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/signin",
		Body:       string(requestBody),
	}

	resp, err := Handler(context.Background(), req)

	if err != nil {
		t.Fatalf("Error calling Lambda handler: %v", err)
	}

	assert.Equal(t, 401, resp.StatusCode, "Expected status code to be 401")
	assert.Contains(t, resp.Body, "Failed to sign in user", "Expected failure message")
}

// パスワードリセットのリクエストが成功することを確認
func TestForgotPasswordHandler_Success(t *testing.T) {
	requestBody, _ := json.Marshal(map[string]string{
		"email": "tensho.arai@ities-inc.co.jp",
	})

	req := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/forgot-password",
		Body:       string(requestBody),
	}

	resp, err := Handler(context.Background(), req)

	if err != nil {
		t.Fatalf("Error calling Lambda handler: %v", err)
	}

	assert.Equal(t, 200, resp.StatusCode, "Expected status code to be 200")
	assert.Contains(t, resp.Body, "Password reset requested", "Expected reset request message")
}
