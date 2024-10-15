# Cognito Lambda Handler (Golang)

このプロジェクトは、Golangで書かれたLambda関数ハンドラーです。Amazon ECRからコンテナイメージを使用して、AWS Lambda上でデプロイすることを想定しています。

## 必要要件

- [Docker](https://www.docker.com/get-started)
- [AWS SAM CLI](https://docs.aws.amazon.com/ja_jp/serverless-application-model/latest/developerguide/install-sam-cli.html)
- Golang 1.19 以上

## セットアップ

1. リポジトリをクローンします:
   ```bash
   git clone https://github.com/your-repo/cognito-lambda-handler-golang.git
   cd cognito-lambda-handler-golang
   ```
2. 依存関係をインストールします:
   ```bash
   go mod tidy
   ```

## ローカル開発とテスト

AWS SAM CLI と Docker を使用して、ローカルでLambda関数を実行することができます。以下の手順に従って、Lambda関数をローカルでビルドし、実行します。

### Dockerイメージのビルド

まず、docker-compose を使用してDockerイメージをビルドします。

```bash
docker-compose build
```

このコマンドは、提供された Dockerfile を使用してDockerイメージをビルドします。

### SAMを使ったローカルでのLambda関数の実行

Dockerイメージがビルドされたら、AWS SAM CLI を使ってローカルでLambda関数を実行できます。

SAMプロジェクトをビルドします:

```bash
sam build --template-file cloudformation/local-template.yaml
```

このコマンドは、Lambda関数をパッケージ化し、ローカル実行用に準備します。

ローカルでLambda関数を実行します:

```bash
sam local invoke MyLambdaFunction --template-file cloudformation/local-template.yaml
```

これにより、ビルドしたDockerイメージを使って、ローカルでLambda関数が実行されます。次のような出力が表示されるはずです:

```bash
"Hello, World"
```


### イベントペイロードを使用したテスト

特定のイベントペイロードでLambda関数をテストするには、`event.json` ファイルを作成し、入力データを指定します。

例: `event.json`

```json
{ "key": "value" }
```

次に、以下のコマンドでイベントを使ってLambda関数を実行します:

```bash
sam local invoke MyLambdaFunction --template-file cloudformation/local-template.yaml --event event.json --env-vars env.json
```

これにより、Lambda関数にイベントが送信された場合の動作をシミュレートできます。

### AWSへのデプロイ

ローカルで関数のテストが完了したら、次のコマンドを使用してAWSにデプロイできます:

```bash
sam deploy --guided
```

このコマンドを実行すると、スタック名、リージョン、IAMロールなどのパラメータを指定しながら、対話形式でデプロイできます。



