package main

import (
	"cognito-lambda-handler/internal/cognito"
	"cognito-lambda-handler/routes"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
	"os"
	"strings"
)

var cognitoService *cognito.Service

// ResponseWriter APIGatewayProxyResponse用のカスタムResponseWriter
type ResponseWriter struct {
	StatusCode int
	Headers    map[string]string
	Body       string
}

// Headerは、ヘッダーのマップを返します
func (rw *ResponseWriter) Header() http.Header {
	return http.Header{}
}

// Write メソッドは、レスポンスボディを記録します
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.Body = string(b)
	return len(b), nil
}

// WriteHeader メソッドは、ステータスコードを設定します
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
}

// NewRequest APIGatewayリクエストをHTTPリクエストに変換
func NewRequest(req events.APIGatewayProxyRequest) *http.Request {
	body := strings.NewReader(req.Body)
	httpReq, _ := http.NewRequest(req.HTTPMethod, req.Path, body)
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}
	return httpReq
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r := routes.RegisterRoutes(cognitoService)

	httpReq := NewRequest(req)

	rw := &ResponseWriter{Headers: map[string]string{}}
	r.ServeHTTP(rw, httpReq)

	return events.APIGatewayProxyResponse{
		StatusCode: rw.StatusCode,
		Headers:    rw.Headers,
		Body:       rw.Body,
	}, nil
}

func init() {
	clientID := os.Getenv("AWS_COGNITO_CLIENT_ID")
	if clientID == "" {
		log.Fatal("AWS_COGNITO_CLIENT_ID is not set")
	}

	var err error
	cognitoService, err = cognito.NewCognitoService(clientID)
	if err != nil {
		log.Fatalf("Failed to initialize Cognito service: %v", err)
	}
}

func main() {
	lambda.Start(handler)
}
