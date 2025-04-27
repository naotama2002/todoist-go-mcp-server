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

- [ ] Go モジュールの初期化
- [ ] 必要なライブラリの追加
  - mcp-go フレームワーク
  - HTTP クライアントライブラリ
  - 設定ファイル管理ライブラリ
- [ ] ディレクトリ構造の設定（github-mcp-server を参考に）

### 2. Todoist API クライアントの実装 (1-2日目)

- [ ] pkg/todoist/client.go の実装
  - [ ] HTTP クライアントの基本実装
  - [ ] 認証ヘッダーの設定
  - [ ] エラーハンドリングの実装
  - [ ] レート制限への対応
  - [ ] タスク関連 API の実装
    - [ ] タスク一覧取得
    - [ ] タスク詳細取得
    - [ ] タスク作成
    - [ ] タスク更新
    - [ ] タスク完了/再開
    - [ ] タスク削除
  - [ ] プロジェクト関連 API の実装
    - [ ] プロジェクト一覧取得
    - [ ] プロジェクト詳細取得
    - [ ] プロジェクト作成
    - [ ] プロジェクト更新
    - [ ] プロジェクト削除

### 3. MCP ツールの実装 (3-4日目)

- [ ] pkg/todoist/tools.go の実装
  - [ ] ツール定義の共通処理
  - [ ] エラーハンドリング

#### タスク管理ツール
- [ ] pkg/todoist/tasks.go の実装
  - [ ] `todoist_get_tasks` の実装
  - [ ] `todoist_get_task` の実装
  - [ ] `todoist_create_task` の実装
  - [ ] `todoist_update_task` の実装
  - [ ] `todoist_close_task` の実装
  - [ ] `todoist_reopen_task` の実装
  - [ ] `todoist_delete_task` の実装

#### プロジェクト管理ツール
- [ ] pkg/todoist/projects.go の実装
  - [ ] `todoist_get_projects` の実装
  - [ ] `todoist_get_project` の実装
  - [ ] `todoist_create_project` の実装
  - [ ] `todoist_update_project` の実装
  - [ ] `todoist_delete_project` の実装

### 4. MCP サーバの実装 (5日目)

- [ ] pkg/todoist/server.go の実装
  - [ ] サーバの初期化
  - [ ] ツールの登録
  - [ ] エラーハンドリングの実装
  - [ ] ログ出力の設定
- [ ] cmd/todoist-mcp-server/main.go の実装
  - [ ] コマンドライン引数の処理
  - [ ] 設定の読み込み
  - [ ] サーバの起動

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
