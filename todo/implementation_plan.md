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
  - [x] プロジェクト関連 API の実装
    - [x] プロジェクト一覧取得
    - [x] プロジェクト詳細取得

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
  - [x] `todoist_delete_task` の実装

#### プロジェクト管理ツール
- [x] pkg/todoist/projects.go の実装
  - [x] `todoist_get_projects` の実装
  - [x] `todoist_get_project` の実装

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

- [x] pkg/toolsets/ の実装
  - [x] デフォルトツールセットの定義
  - [x] カスタムツールセットの実装（必要に応じて）

### 6. テストの実装 (6日目)

- [x] ユニットテストの実装
  - [x] client_test.go
  - [x] tasks_test.go
  - [x] projects_test.go
  - [x] server_test.go

### 7. ドキュメントの整備と最終調整 (7日目)

- [x] README.md の作成
  - https://github.com/github/github-mcp-server の実装を参考にしてください。/tmp/github-mcp-server に clone されています。
  - MCP 利用ユーザが迷わない記述を心がけてください。
- [x] API ドキュメントの更新
- [x] 使用例の追加
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
1. プロジェクト構造の設定
   - 基本的なディレクトリ構造の作成
   - Go モジュールの初期化
   - 必要なパッケージの追加

2. Todoist API クライアントの実装
   - HTTP クライアントの実装
   - 認証機能の実装
   - タスク関連 API の実装
   - プロジェクト関連 API の実装

3. MCP サーバーの実装
   - サーバー構造の実装
   - HTTP モードと stdio モードのサポート
   - ツールの登録機能

4. MCP ツールの実装
   - タスク管理ツール
     - `todoist_get_tasks`
     - `todoist_get_task`
     - `todoist_create_task`
     - `todoist_update_task`
     - `todoist_close_task`
     - `todoist_delete_task`
   - プロジェクト管理ツール
     - `todoist_get_projects`
     - `todoist_get_project`

5. ツールセットの実装
   - デフォルトツールセットの定義
   - カスタムツールセットの実装

6. ユニットテストの実装
   - client_test.go
   - tasks_test.go
   - projects_test.go
   - server_test.go

7. ドキュメントの整備
   - README.md の作成
   - API ドキュメントの更新（docs/api.md）
   - 使用例の追加（docs/examples.md）

### 残りのタスク

1. コードのリファクタリングと最適化
