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
or
```bash
sam.cmd build --template-file cloudformation/local-template.yaml
```

このコマンドは、Lambda関数をパッケージ化し、ローカル実行用に準備します。

ローカルでLambda関数を実行します:

```bash
sam local invoke MyLambdaFunction --template-file cloudformation/local-template.yaml
```
or
```bash
sam.cmd local invoke MyLambdaFunction --template-file cloudformation/local-template.yaml
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
or
```bash
sam.cmd local invoke MyLambdaFunction --template-file cloudformation/local-template.yaml --event event.json --env-vars env.json
```

これにより、Lambda関数にイベントが送信された場合の動作をシミュレートできます。

### AWSへのデプロイ

ローカルで関数のテストが完了したら、次のコマンドを使用してAWSにデプロイできます:

```bash
sam deploy --guided
```

このコマンドを実行すると、スタック名、リージョン、IAMロールなどのパラメータを指定しながら、対話形式でデプロイできます。


## ローカルでAPIサーバ立ち上げてリクエスト

以下の手順に従って、サーバーをローカルで起動し、Cognitoのサインアップおよびサインインをテストします。

### サーバーのビルドと起動

1. Docker Compose を使用してプロジェクトをビルドします:

   ```bash
   docker compose build
   ```

2. SAM CLI を使用してプロジェクトをビルドします:

   ```bash
   sam.cmd build --template-file cloudformation/local-template.yaml
   ```

3. SAM CLI を使用してローカルAPIを起動します:

   ```bash
   sam.cmd local start-api --template-file cloudformation/local-template.yaml --env-vars env.json
   ```

   上記のコマンドにより、ローカルサーバーが `http://127.0.0.1:3000` で起動します。

---

### サインアップのリクエスト

以下のコマンドを使用して、ユーザーをCognitoにサインアップします:

```bash
curl -X POST http://127.0.0.1:3000/signup -H "Content-Type: application/json" -d '{"email": "testuser@example.com", "password": "Password123!", "phone_number": "+1234567890", "given_name": "Test", "family_name": "User"}'
```

リクエストボディには以下の情報を含めます:
- **email**: ユーザーのメールアドレス
- **password**: ユーザーのパスワード
- **phone_number**: ユーザーの電話番号
- **given_name**: ユーザーの名
- **family_name**: ユーザーの姓

---

### サインインのリクエスト

サインアップ後、以下のコマンドを使用してサインインを行います:

```bash
curl -X POST http://127.0.0.1:3000/signin -H "Content-Type: application/json" -d '{"email": "test@example.com", "password": "Password123!"}'
```

リクエストボディには以下の情報を含めます:
- **email**: サインインするユーザーのメールアドレス
- **password**: ユーザーのパスワード

サインインに成功すると、Cognitoからアクセストークンが返されます。このトークンは、認証が必要なリソースにアクセスするために使用できます。

---