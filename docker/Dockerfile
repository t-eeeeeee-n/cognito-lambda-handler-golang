# Golangベースのイメージを使用
FROM golang:1.23 AS builder

# 作業ディレクトリの作成
WORKDIR /go/src/app

# ソースコードをコンテナにコピー
COPY . /go/src/app/

# 依存関係をダウンロード
RUN go mod tidy

# Goアプリケーションを静的にビルド
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main cmd/lambda_handler/main.go

# Lambdaの実行環境として公式のAWS Lambdaベースイメージを使用
FROM public.ecr.aws/lambda/go:1

# ビルド済みのバイナリをコピー
COPY --from=builder /main ${LAMBDA_TASK_ROOT}

# デフォルトのLambdaエントリポイントを指定
CMD [ "main" ]
