package main

import (
	"cognito-lambda-handler/internal/cognito"
	"cognito-lambda-handler/routes"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
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

// Header ヘッダーのマップを返します
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

func Handler(_ context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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
	var err error
	if _, exists := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME"); !exists {
		// ローカル環境でのみ .env をロード
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	clientId := os.Getenv("AWS_COGNITO_CLIENT_ID")
	if clientId == "" {
		log.Fatal("AWS_COGNITO_CLIENT_ID is not set")
	}

	clientSecret := os.Getenv("AWS_COGNITO_CLIENT_SECRET")
	if clientSecret == "" {
		log.Fatal("AWS_COGNITO_CLIENT_SECRET is not set")
	}

	poolId := os.Getenv("AWS_COGNITO_POOL_ID")
	if poolId == "" {
		log.Fatal("AWS_COGNITO_POOL_ID is not set")
	}

	cognitoService, err = cognito.NewCognitoService(clientId, clientSecret, poolId)
	if err != nil {
		log.Fatalf("Failed to initialize Cognito service: %v", err)
	}
}

func main() {
	if _, isLambda := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME"); isLambda {
		// AWS Lambda環境
		lambda.Start(Handler)
	} else {
		// ローカル環境
		r := routes.RegisterRoutes(cognitoService)
		log.Println("Starting local server on :8080")
		log.Fatal(http.ListenAndServe(":8080", r))
	}
}
