# Todoist MCP Server 仕様書

## 概要

Todoist MCP Server は、Todoist API v1 を利用して、タスクとプロジェクトの管理機能を MCP ツールとして提供します。このサーバは Go 言語と公式 MCP Go SDK (modelcontextprotocol/go-sdk) を使用して実装されます。

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

#### 1. `todoist_get_task_filter_rules`

Todoistタスクフィルターのルールと例を取得します。このツールは、自然言語クエリをTodoistフィルター構文に変換するために使用されます。

**パラメータ:**
なし

**戻り値:**
Markdown形式のテキストで、以下の情報を含みます：

```
# Introduction to Filters

Filters in Todoist are custom views that display tasks based on specific criteria. You can filter tasks by name, date, project, label, priority, creation date, and more.

## Filter Symbols

| Symbol | Meaning | Example |
|--------|---------|--------|
| \| | OR | today \| overdue |
| & | AND | today & p1 |
| ! | NOT | !subtask |
| () | Priority processing | (today \| overdue) & #work |
| , | Display in separate lists | date: yesterday, today |
| \\ | Use special characters as regular characters | #1 \\& 2 |

## Advanced Queries

### Keyword-Based Filters
| Description | Query |
|-------------|-------|
| Tasks containing "meeting" | search: meeting |
| Tasks containing "meeting" scheduled for today | search: meeting & today |

...

## Useful Filter Examples

| Description | Query |
|-------------|-------|
| Overdue or today's tasks in "Work" project | (today \| overdue) & #work |
| Tasks with no date | no date |
```

完全なフィルター構文リファレンスには、基本フィルター、論理演算子、日付フィルター、優先度フィルター、プロジェクトとラベルのフィルター、検索フィルター、その他の便利なフィルター例が含まれます。
```

#### 2. `todoist_get_tasks`

タスクの一覧を取得します。

**パラメータ:**
- `projectId` (string, optional): プロジェクトID
- `filter` (string, optional): フィルター文字列

**戻り値:**
```json
{
  "tasks": [
    {
      "id": "2995104339",
      "content": "Buy Milk",
      "description": "",
      "project_id": "2203306141",
      "parent_id": null,
      "priority": 1,
      "due": {
        "date": "2016-09-01",
        "is_recurring": false,
        "datetime": "2016-09-01T12:00:00.000000Z",
        "string": "tomorrow at 12",
        "timezone": "Europe/Moscow"
      }
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
    "parent_id": null,
    "priority": 1,
    "due": {
      "date": "2016-09-01",
      "is_recurring": false,
      "datetime": "2016-09-01T12:00:00.000000Z",
      "string": "tomorrow at 12",
      "timezone": "Europe/Moscow"
    }
  }
}
```

#### 3. `todoist_create_task`

新しいタスクを作成します。

**パラメータ:**
- `content` (string, required): タスクの内容
- `description` (string, optional): タスクの詳細説明
- `projectId` (string, optional): プロジェクトID
- `parentId` (string, optional): 親タスクID
- `priority` (int, optional): 優先度 (1-4)
- `dueString` (string, optional): 期限を表す文字列
- `dueDate` (string, optional): 期限の日付
- `dueDatetime` (string, optional): 期限の日時

**戻り値:**
```json
{
  "task": {
    "id": "2995104339",
    "content": "Buy Milk",
    "description": "",
    "project_id": "2203306141",
    "parent_id": null,
    "priority": 1,
    "due": {
      "date": "2016-09-01",
      "is_recurring": false,
      "datetime": "2016-09-01T12:00:00.000000Z",
      "string": "tomorrow at 12",
      "timezone": "Europe/Moscow"
    }
  }
}
```

#### 4. `todoist_update_task`

既存のタスクを更新します。

**パラメータ:**
- `id` (string, required): タスクID
- `content` (string, optional): タスクの内容
- `description` (string, optional): タスクの詳細説明
- `priority` (int, optional): 優先度 (1-4)
- `dueString` (string, optional): 期限を表す文字列
- `dueDate` (string, optional): 期限の日付
- `dueDatetime` (string, optional): 期限の日時

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

#### 6. `todoist_delete_task`

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
      "child_order": 1,
      "color": "charcoal",
      "is_shared": false,
      "is_favorite": false,
      "inbox_project": false,
      "view_style": "list"
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
    "child_order": 1,
    "color": "charcoal",
    "is_shared": false,
    "is_favorite": false,
    "inbox_project": false,
    "view_style": "list"
  }
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