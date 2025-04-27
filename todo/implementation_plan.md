# Todoist MCP サーバ実装計画

## 概要
Todoist REST API v2 を利用した MCP サーバを Go 言語で実装します。このサーバは、タスクとプロジェクトの管理機能を MCP ツールとして提供します。

## プロジェクト構造
GitHub の github-mcp-server を参考にした構造で実装します。

```
/
├── cmd/
│   └── todoist-mcp-server/
│       └── main.go       # サーバのエントリーポイント
├── pkg/
│   ├── todoist/          # Todoist 関連の実装
│   │   ├── server.go     # サーバー実装
│   │   ├── tools.go      # ツール定義
│   │   ├── tasks.go      # タスク関連ツール
│   │   ├── projects.go   # プロジェクト関連ツール
│   │   └── client.go     # Todoist API クライアント
│   ├── log/              # ログ関連
│   └── toolsets/         # ツールセット定義
├── docs/                 # ドキュメント
│   └── spec.md           # 仕様書
├── script/               # スクリプト
└── go.mod, go.sum        # 依存関係ファイル
```

## 実装ステップ

### 1. プロジェクト初期化と依存関係の設定 (1日目)

- [x] Go モジュールの初期化
- [x] 必要なライブラリの追加
  - [x] mcp-go フレームワーク
  - [x] HTTP クライアントライブラリ
  - [x] 設定ファイル管理ライブラリ
- [x] ディレクトリ構造の設定（github-mcp-server を参考に）

### 2. Todoist API クライアントの実装 (1-2日目)

- [x] pkg/todoist/client.go の実装
  - [x] HTTP クライアントの基本実装
  - [x] 認証ヘッダーの設定
  - [x] エラーハンドリングの実装
  - [x] レート制限への対応
  - [x] タスク関連 API の実装
    - [x] タスク一覧取得
    - [x] タスク詳細取得
    - [x] タスク作成
    - [x] タスク更新
    - [x] タスク完了/再開
    - [x] タスク削除
  - [ ] プロジェクト関連 API の実装
    - [ ] プロジェクト一覧取得
    - [ ] プロジェクト詳細取得

### 3. MCP ツールの実装 (3-4日目)

- [x] pkg/todoist/tools.go の実装
  - [x] ツール定義の共通処理
  - [x] エラーハンドリング

#### タスク管理ツール
- [x] pkg/todoist/tasks.go の実装
  - [x] `todoist_get_tasks` の実装
  - [x] `todoist_get_task` の実装
  - [x] `todoist_create_task` の実装
  - [x] `todoist_update_task` の実装
  - [x] `todoist_close_task` の実装
  - [x] `todoist_reopen_task` の実装
  - [x] `todoist_delete_task` の実装

#### プロジェクト管理ツール
- [ ] pkg/todoist/projects.go の実装
  - [ ] `todoist_get_projects` の実装
  - [ ] `todoist_get_project` の実装

### 4. MCP サーバの実装 (5日目)

- [x] pkg/todoist/server.go の実装
  - [x] サーバの初期化
  - [x] ツールの登録
  - [x] エラーハンドリングの実装
  - [x] ログ出力の設定
- [x] cmd/todoist-mcp-server/main.go の実装
  - [x] コマンドライン引数の処理
  - [x] 設定の読み込み
  - [x] サーバの起動

### 5. ツールセットの実装 (5日目)

- [ ] pkg/toolsets/ の実装
  - [ ] デフォルトツールセットの定義
  - [ ] カスタムツールセットの実装（必要に応じて）

### 6. テストの実装 (6日目)

- [ ] ユニットテストの実装
  - [ ] client_test.go
  - [ ] tasks_test.go
  - [ ] projects_test.go
  - [ ] server_test.go
- [ ] 統合テストの実装
- [ ] モックサーバの実装（必要に応じて）

### 7. ドキュメントの整備と最終調整 (7日目)

- [ ] README.md の作成
- [ ] API ドキュメントの更新
- [ ] 使用例の追加
- [ ] コードのリファクタリングと最適化

## 実装の注意点

### エラーハンドリング
- Todoist API からのエラーレスポンスを適切に処理
- ネットワークエラーの処理
- API レート制限への対応

### セキュリティ
- API トークンの安全な管理
- 環境変数または暗号化された設定ファイルでの保存

### パフォーマンス
- 適切なキャッシュ戦略の検討
- 並行処理の活用（必要に応じて）

## 技術スタック
- 言語: Go
- フレームワーク: mcp-go (https://github.com/mark3labs/mcp-go)
- API: Todoist REST API v2 (https://developer.todoist.com/rest/v2/)

## 参考リソース
- GitHub MCP サーバ実装: https://github.com/github/github-mcp-server
- Todoist API ドキュメント: https://developer.todoist.com/rest/v2/
- mcp-go ドキュメント: https://github.com/mark3labs/mcp-go

## 実装状況（2025-04-27 更新）

### 完了した作業
1. プロジェクト初期化と依存関係の設定
   - Go モジュールの初期化
   - 必要なライブラリの追加
   - ディレクトリ構造の設定

2. Todoist API クライアントの基本実装
   - HTTP クライアントの実装
   - 認証ヘッダーの設定
   - エラーハンドリングの実装
   - タスク関連 API の実装
     - タスク一覧取得
     - タスク詳細取得
     - タスク作成
     - タスク更新
     - タスク完了/再開
     - タスク削除

3. MCP ツールの実装
   - ツール定義の共通処理
   - タスク管理ツールの実装
     - `todoist_get_tasks` ツールの実装
     - `todoist_get_task` ツールの実装
     - `todoist_create_task` ツールの実装
     - `todoist_update_task` ツールの実装
     - `todoist_close_task` ツールの実装
     - `todoist_reopen_task` ツールの実装
     - `todoist_delete_task` ツールの実装

4. MCP サーバの実装
   - サーバの初期化
   - ツールの登録
   - HTTP モードと stdio モードの両方をサポート
   - MCP プロトコルに準拠したリクエスト処理の実装

### 現在の動作確認状況
- HTTP モードでのサーバ起動が可能
- stdio モードでのサーバ起動が可能
- すべてのタスク管理ツールが実装済み
- MCP プロトコルに準拠したリクエスト処理が実装済み
- Claude Desktop などの MCP クライアントから呼び出し可能
- テスト用の MCP クライアントを実装し、コマンドライン引数からツール名と引数を指定可能

### 次のステップ
1. プロジェクト管理ツールの実装
   - `todoist_get_projects`
   - `todoist_get_project`

2. ユニットテストの追加
3. ドキュメントの整備
