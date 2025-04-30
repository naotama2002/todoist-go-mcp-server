# Todoist MCP サーバーテスト計画

## 概要
Todoist MCP サーバーの品質を確保するために、以下のテスト戦略を実施します。

## テスト対象ファイル

### 1. クライアントテスト (`client_test.go`)
Todoist API クライアントの機能をテストします。

- **テスト対象**: `pkg/todoist/client.go`
- **テスト内容**:
  - クライアント初期化のテスト
  - API リクエスト処理のテスト
  - エラーハンドリングのテスト
  - 各APIメソッドのテスト
    - GetTasks
    - GetTask
    - GetProjects
    - GetProject
    - CreateTask
    - UpdateTask
    - CloseTask
    - ReopenTask
    - DeleteTask

### 2. タスク関連ツールテスト (`tasks_test.go`)
タスク関連のMCPツール実装をテストします。

- **テスト対象**: `pkg/todoist/tasks.go`
- **テスト内容**:
  - ツール定義のテスト
    - スキーマの検証
    - 必須パラメータの検証
  - ハンドラ関数のテスト
    - HandleGetTasks
    - HandleGetTask
    - HandleCreateTask
    - HandleUpdateTask
    - HandleCloseTask
    - HandleDeleteTask
  - パラメータ処理関数のテスト
    - OptionalParam
    - RequiredParam
    - OptionalStringArrayParam

### 3. プロジェクト関連ツールテスト (`projects_test.go`)
プロジェクト関連のMCPツール実装をテストします。

- **テスト対象**: `pkg/todoist/projects.go`
- **テスト内容**:
  - ツール定義のテスト
    - スキーマの検証
  - ハンドラ関数のテスト
    - HandleGetProjects
    - HandleGetProject

### 4. サーバーテスト (`server_test.go`)
MCP サーバー実装をテストします。

- **テスト対象**: `pkg/todoist/server.go`
- **テスト内容**:
  - サーバー初期化のテスト
  - ツール登録のテスト
  - HTTP リクエスト処理のテスト
  - stdio モードのテスト

## テスト戦略

### モック戦略
実際の Todoist API を呼び出さないようにするため、以下のモックアプローチを採用します：

1. **HTTPモック**: `httptest` パッケージを使用してHTTPリクエストをモック
2. **インターフェースモック**: クライアントインターフェースを定義し、テスト用モックを実装

### テストカバレッジ目標
- コード全体で 80% 以上のカバレッジを目指す
- 重要なエラーハンドリングパスのカバレッジを確保

### テスト実行方法
```bash
# 全テストの実行
go test ./...

# カバレッジレポートの生成
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## 実装順序

1. モックの実装
2. クライアントテスト (`client_test.go`)
3. タスク関連ツールテスト (`tasks_test.go`)
4. プロジェクト関連ツールテスト (`projects_test.go`)
5. サーバーテスト (`server_test.go`)
