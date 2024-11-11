## SQLCプロジェクトの準備
- tools.goの作成
- go mod tidy の実行

- go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

- go run github.com/sqlc-dev/sqlc/cmd/sqlc init 

## データベースの構築
- makefileに定義が書いてあるはず

## コードの自動生成
- go run github.com/sqlc-dev/sqlc/cmd/sqlc generateを実行