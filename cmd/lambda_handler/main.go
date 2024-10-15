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
)

var cognitoService *cognito.Service

// ResponseWriter カスタムResponseWriterの定義
type ResponseWriter struct {
	StatusCode int
	Headers    http.Header
	Body       string
}

func (rw *ResponseWriter) Header() http.Header {
	if rw.Headers == nil {
		rw.Headers = http.Header{}
	}
	return rw.Headers
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.Body = string(b)
	return len(b), nil
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
}

func (rw *ResponseWriter) ToAPIGatewayProxyResponse() events.APIGatewayProxyResponse {
	headers := make(map[string]string)
	for k, v := range rw.Headers {
		headers[k] = v[0]
	}
	return events.APIGatewayProxyResponse{
		StatusCode: rw.StatusCode,
		Headers:    headers,
		Body:       rw.Body,
	}
}

// NewRequest APIGatewayリクエストをHTTPリクエストに変換
func NewRequest(req events.APIGatewayProxyRequest) *http.Request {
	httpReq, _ := http.NewRequest(req.HTTPMethod, req.Path, nil)
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}
	return httpReq
}

// Lambdaハンドラとしてリクエストを処理
func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// ルーティングを登録
	r := routes.RegisterRoutes(cognitoService)

	// カスタムResponseWriterでレスポンスをキャプチャ
	rw := &ResponseWriter{}
	r.ServeHTTP(rw, NewRequest(req))

	// APIGatewayProxyResponseとしてレスポンスを返す
	return rw.ToAPIGatewayProxyResponse(), nil
}

func init() {
	// 環境変数からCognitoのクライアントIDを取得
	clientID := os.Getenv("AWS_COGNITO_CLIENT_ID")
	if clientID == "" {
		log.Fatal("AWS_COGNITO_CLIENT_ID is not set")
	}

	// Cognitoサービスの初期化
	var err error
	cognitoService, err = cognito.NewCognitoService(clientID)
	if err != nil {
		log.Fatalf("Failed to initialize Cognito service: %v", err)
	}
}

func main() {
	// Lambda関数のハンドラを起動
	lambda.Start(handler)
}
