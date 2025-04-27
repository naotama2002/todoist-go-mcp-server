# Todoist MCP Server 仕様書

## 概要

Todoist MCP Server は、Todoist REST API v2 を利用して、タスクとプロジェクトの管理機能を MCP ツールとして提供します。このサーバは Go 言語と mcp-go フレームワークを使用して実装されます。

## 認証

Todoist API を利用するためには、Todoist の個人用 API トークンが必要です。このトークンは環境変数 `TODOIST_API_TOKEN` として設定するか、設定ファイルで指定します。

```go
// 認証トークンの取得例
token := os.Getenv("TODOIST_API_TOKEN")
if token == "" {
    // 設定ファイルからの読み込みなど、代替手段を実装
}
```

## MCP ツール

### タスク管理

#### 1. `todoist_get_tasks`

タスクの一覧を取得します。

**パラメータ:**
- `projectId` (string, optional): プロジェクトID
- `sectionId` (string, optional): セクションID
- `label` (string, optional): ラベル名
- `filter` (string, optional): フィルター文字列
- `lang` (string, optional): 言語コード
- `ids` ([]string, optional): タスクIDのリスト

**戻り値:**
```json
{
  "tasks": [
    {
      "id": "2995104339",
      "content": "Buy Milk",
      "description": "",
      "project_id": "2203306141",
      "section_id": "7025",
      "parent_id": "2995104589",
      "labels": ["Food", "Shopping"],
      "priority": 1,
      "due": {
        "date": "2016-09-01",
        "is_recurring": false,
        "datetime": "2016-09-01T12:00:00.000000Z",
        "string": "tomorrow at 12",
        "timezone": "Europe/Moscow"
      },
      "url": "https://todoist.com/showTask?id=2995104339"
    }
  ]
}
```

#### 2. `todoist_get_task`

指定されたIDのタスクを取得します。

**パラメータ:**
- `id` (string, required): タスクID

**戻り値:**
```json
{
  "task": {
    "id": "2995104339",
    "content": "Buy Milk",
    "description": "",
    "project_id": "2203306141",
    "section_id": "7025",
    "parent_id": "2995104589",
    "labels": ["Food", "Shopping"],
    "priority": 1,
    "due": {
      "date": "2016-09-01",
      "is_recurring": false,
      "datetime": "2016-09-01T12:00:00.000000Z",
      "string": "tomorrow at 12",
      "timezone": "Europe/Moscow"
    },
    "url": "https://todoist.com/showTask?id=2995104339"
  }
}
```

#### 3. `todoist_create_task`

新しいタスクを作成します。

**パラメータ:**
- `content` (string, required): タスクの内容
- `description` (string, optional): タスクの詳細説明
- `projectId` (string, optional): プロジェクトID
- `sectionId` (string, optional): セクションID
- `parentId` (string, optional): 親タスクID
- `labels` ([]string, optional): ラベルのリスト
- `priority` (int, optional): 優先度 (1-4)
- `dueString` (string, optional): 期限を表す文字列
- `dueDate` (string, optional): 期限の日付
- `dueDatetime` (string, optional): 期限の日時
- `dueLang` (string, optional): 期限の言語

**戻り値:**
```json
{
  "task": {
    "id": "2995104339",
    "content": "Buy Milk",
    "description": "",
    "project_id": "2203306141",
    "section_id": "7025",
    "parent_id": "2995104589",
    "labels": ["Food", "Shopping"],
    "priority": 1,
    "due": {
      "date": "2016-09-01",
      "is_recurring": false,
      "datetime": "2016-09-01T12:00:00.000000Z",
      "string": "tomorrow at 12",
      "timezone": "Europe/Moscow"
    },
    "url": "https://todoist.com/showTask?id=2995104339"
  }
}
```

#### 4. `todoist_update_task`

既存のタスクを更新します。

**パラメータ:**
- `id` (string, required): タスクID
- `content` (string, optional): タスクの内容
- `description` (string, optional): タスクの詳細説明
- `labels` ([]string, optional): ラベルのリスト
- `priority` (int, optional): 優先度 (1-4)
- `dueString` (string, optional): 期限を表す文字列
- `dueDate` (string, optional): 期限の日付
- `dueDatetime` (string, optional): 期限の日時
- `dueLang` (string, optional): 期限の言語

**戻り値:**
```json
{
  "success": true
}
```

#### 5. `todoist_close_task`

タスクを完了状態にします。

**パラメータ:**
- `id` (string, required): タスクID

**戻り値:**
```json
{
  "success": true
}
```

#### 6. `todoist_reopen_task`

完了したタスクを再開します。

**パラメータ:**
- `id` (string, required): タスクID

**戻り値:**
```json
{
  "success": true
}
```

#### 7. `todoist_delete_task`

タスクを削除します。

**パラメータ:**
- `id` (string, required): タスクID

**戻り値:**
```json
{
  "success": true
}
```

### プロジェクト管理

#### 1. `todoist_get_projects`

プロジェクトの一覧を取得します。

**パラメータ:**
なし

**戻り値:**
```json
{
  "projects": [
    {
      "id": "2203306141",
      "name": "Shopping List",
      "comment_count": 10,
      "order": 1,
      "color": "charcoal",
      "is_shared": false,
      "is_favorite": false,
      "is_inbox_project": false,
      "is_team_inbox": false,
      "view_style": "list",
      "url": "https://todoist.com/showProject?id=2203306141"
    }
  ]
}
```

#### 2. `todoist_get_project`

指定されたIDのプロジェクトを取得します。

**パラメータ:**
- `id` (string, required): プロジェクトID

**戻り値:**
```json
{
  "project": {
    "id": "2203306141",
    "name": "Shopping List",
    "comment_count": 10,
    "order": 1,
    "color": "charcoal",
    "is_shared": false,
    "is_favorite": false,
    "is_inbox_project": false,
    "is_team_inbox": false,
    "view_style": "list",
    "url": "https://todoist.com/showProject?id=2203306141"
  }
}
```

#### 3. `todoist_create_project`

新しいプロジェクトを作成します。

**パラメータ:**
- `name` (string, required): プロジェクト名
- `parentId` (string, optional): 親プロジェクトID
- `color` (string, optional): 色
- `isFavorite` (bool, optional): お気に入りかどうか
- `viewStyle` (string, optional): 表示スタイル

**戻り値:**
```json
{
  "project": {
    "id": "2203306141",
    "name": "Shopping List",
    "comment_count": 0,
    "order": 1,
    "color": "charcoal",
    "is_shared": false,
    "is_favorite": false,
    "is_inbox_project": false,
    "is_team_inbox": false,
    "view_style": "list",
    "url": "https://todoist.com/showProject?id=2203306141"
  }
}
```

#### 4. `todoist_update_project`

既存のプロジェクトを更新します。

**パラメータ:**
- `id` (string, required): プロジェクトID
- `name` (string, optional): プロジェクト名
- `color` (string, optional): 色
- `isFavorite` (bool, optional): お気に入りかどうか
- `viewStyle` (string, optional): 表示スタイル

**戻り値:**
```json
{
  "success": true
}
```

#### 5. `todoist_delete_project`

プロジェクトを削除します。

**パラメータ:**
- `id` (string, required): プロジェクトID

**戻り値:**
```json
{
  "success": true
}
```

## 実装計画

1. プロジェクト構造の設定
2. 基本的な MCP サーバの実装
3. Todoist API クライアントの実装
4. 各 MCP ツールの実装
5. エラーハンドリングの実装
6. テストの実装
7. ドキュメントの整備

## エラーハンドリング

Todoist API からのエラーレスポンスは、適切なエラーメッセージとステータスコードで返します。

```json
{
  "error": "Invalid request",
  "code": 400
}
```

## API レート制限

Todoist API にはレート制限があります。制限に達した場合は、適切なエラーメッセージを返します。

```json
{
  "error": "Rate limit exceeded",
  "code": 429
}