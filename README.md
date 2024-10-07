# GitHub Actionsのテンプレート

このテンプレートは、GitHub Actionsを使用してAWS Elastic Container Registry（ECR）にDockerイメージをプッシュし、その後AWS CloudFormationを使用してデプロイするためのものです。

## 前提条件

- OpenID Connect (OIDC) を使用してAWSに認証することができます。
- 参照資料: [GitHub Actionsを使ったOIDCによるAWS認証の方法](https://zenn.dev/kou_pg_0131/articles/gh-actions-oidc-aws)

## 設定手順

1. **GitHubのSecretに環境変数の情報を登録**
    - `AWS_ACCOUNT_ID`: AWSアカウントID（例：123456789012）
    - `AWS_REGION`: AWSリージョン（例：ap-northeast-1）
    - `AWS_ROLE_TO_ASSUME`: AWS IAMロールのARN（例：arn:aws:iam::123456789012:role/YourRoleName）
    - `AWS_ACCESS_KEY_ID`: [アクセスキーID]
    - `AWS_SECRET_ACCESS_KEY`: [シークレットアクセスキー]

## 使用方法

1. **GitHub Actionsの実行**:
   - リポジトリに変更をコミットし、`main` ブランチにプッシュするとGitHub Actionsのワークフローが自動的に開始されます。

## GitHub Actionsワークフローの概要

ワークフローは以下のステップで構成されています：

1. **リポジトリ名の抽出**:
   - GitHubリポジトリ名を抽出し、後続のステップで使用します。

2. **ECRリポジトリの存在確認と作成**:
   - 指定された名前のECRリポジトリが存在するか確認し、存在しない場合はGitHub Actions内で自動的に作成します。

3. **Dockerイメージのビルドとプッシュ**:
   - AWSの認証情報を設定し、ECRにログインします。
   - Dockerイメージをビルドし、ECRにプッシュします。

4. **AWS CloudFormationによるデプロイ**:
   - CloudFormationテンプレートを使用して、新しいイメージを使用するためのAWSリソースをデプロイまたは更新します。

このテンプレートを使用することで、ソースコードの変更を簡単かつ迅速にAWSにデプロイすることができます。