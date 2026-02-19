package todoist

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetTaskFilterRules はtodoist_get_task_filter_rulesツールを返します
func (tp *ToolProvider) GetTaskFilterRules() mcp.Tool {
	// 入力スキーマを定義
	inputSchema := map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}

	// 入力スキーマをJSONに変換
	inputSchemaJSON, err := json.Marshal(inputSchema)
	if err != nil {
		tp.logger.WithError(err).Error("Failed to marshal input schema")
		return mcp.Tool{}
	}

	return mcp.Tool{
		Name:        "todoist_get_task_filter_rules",
		Description: "Get the filter rules and examples for Todoist task filters. Use this information to translate natural language queries into Todoist filter syntax for the todoist_get_tasks tool.",
		InputSchema: json.RawMessage(inputSchemaJSON),
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}
}

// HandleGetTaskFilterRules はtodoist_get_task_filter_rulesツールリクエストを処理します
func (tp *ToolProvider) HandleGetTaskFilterRules(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// リクエストをログに記録
	tp.logger.Info("Getting task filter rules")

	// フィルタールールのテキストを返却
	filterRulesText := `# Introduction to Filters

Filters in Todoist are custom views that display tasks based on specific criteria. You can filter tasks by name, date, project, label, priority, creation date, and more.

## Creating Filters

1. Select "Filters & Labels" in the sidebar
2. Click the add icon next to Filters
3. Enter a name for your filter and change its color (optional)
4. Enter your filter query
5. Click "Add" to save your filter

## Filter Symbols

| Symbol | Meaning | Example |
|--------|---------|---------|
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
| Tasks containing either "meeting" or "work" | search: meeting \| search: work |
| Tasks containing web links | search: http |

### Subtask Filters

| Description | Query |
|-------------|-------|
| Show all subtasks | subtask |
| Show only parent tasks | !subtask |

### Date-Based Filters

| Description | Query |
|-------------|-------|
| Tasks for a specific date | date: jan 3 |
| Tasks before a specific date | before: may 5 |
| Tasks after a specific date | after: may 5 |
| Overdue and tasks due in the next 4 hours | before: +4 hours |
| Weekday tasks for this week | before: saturday |
| Tasks for next week | (date: next week \| after: next week) & before: 1 week after next week |
| Tasks with no date | no date |
| Tasks with a date | !no date |
| Tasks with date and time | !no date & !no time |
| Overdue tasks | overdue |

### Priority Filters

| Description | Query |
|-------------|-------|
| Priority 1 tasks | p1 |
| Priority 2 tasks | p2 |
| Priority 3 tasks | p3 |
| No priority (priority 4) tasks | no priority |

### Label Filters

| Description | Query |
|-------------|-------|
| Tasks with "email" label | @email |
| Tasks with no labels | no labels |

### Project and Section Filters

| Description | Query |
|-------------|-------|
| Tasks in "Work" project | #work |
| Tasks in "Work" project and its subprojects | ##work |
| Tasks in sections named "Meetings" | /meetings |
| Tasks in "Meetings" section in "Work" project | #work & /meetings |
| Tasks without sections | !/* |

### Creation Date Filters

| Description | Query |
|-------------|-------|
| Tasks created on a specific date | created: jan 3 2023 |
| Tasks created more than 365 days ago | created before: -365 days |
| Tasks created today | created: today |

### Shared Task and Assignment Filters

| Description | Query |
|-------------|-------|
| Tasks assigned to others | assigned to: others |
| Tasks assigned by a specific person | assigned by: [name] |
| Tasks you assigned to others | assigned by: me |
| Tasks that are assigned | assigned |
| Tasks in shared projects | shared |

## Useful Filter Examples

| Description | Query |
|-------------|-------|
| Overdue or today's tasks in "Work" project | (today \| overdue) & #work |
| Tasks with no date | no date |
| Tasks with no time | no time |
| Tasks with @waiting label in the next 7 days | 7 days & @waiting |
| Tasks created more than 30 days ago | created before: -30 days |
| Saturday tasks with the "evening" label | saturday & @evening |
| Tasks assigned to you in the "Work" project | #work & assigned to: me |

## Using Wildcards

Use asterisk (*) to filter tasks with similar strings:

| Description | Query |
|-------------|-------|
| Tasks with labels starting with "home" | @home* |
| Tasks in projects ending with "work" | #*work |
| Tasks in sections containing "meeting" | /*meeting* |

## Running Multiple Filter Queries

Use comma (,) to display multiple lists in the same view:

Example: "p1 & overdue, p4 & today" - Shows priority 1 overdue tasks and priority 4 tasks due today`

	return newToolResultText(filterRulesText), nil
}
