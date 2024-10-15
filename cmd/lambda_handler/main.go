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
	var body *strings.Reader
	if req.Body != "" {
		body = strings.NewReader(req.Body)
	} else {
		body = strings.NewReader("")
	}
	httpReq, _ := http.NewRequest(req.HTTPMethod, req.Path, body)
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}
	return httpReq
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r := routes.RegisterRoutes(cognitoService)
	rw := &ResponseWriter{}
	r.ServeHTTP(rw, NewRequest(req))

	// http.Header (map[string][]string) を map[string]string に変換
	headers := make(map[string]string)
	for k, v := range rw.Headers {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	// エラー時のレスポンス処理
	if rw.StatusCode >= 400 {
		return events.APIGatewayProxyResponse{
			StatusCode: rw.StatusCode,
			Body:       rw.Body,
			Headers:    headers,
		}, nil
	}

	return rw.ToAPIGatewayProxyResponse(), nil
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
