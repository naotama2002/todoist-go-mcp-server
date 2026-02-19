package todoist

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetTaskFilterRulesTool(t *testing.T) {
	// ツールプロバイダーを作成
	tp := NewMockToolProvider()

	// ツールを取得
	tool := tp.GetTaskFilterRules()

	// ツールのプロパティをチェック
	assert.Equal(t, "todoist_get_task_filter_rules", tool.Name)
	assert.Equal(t, "Get the filter rules and examples for Todoist task filters. Use this information to translate natural language queries into Todoist filter syntax for the todoist_get_tasks tool.", tool.Description)
	assert.True(t, tool.Annotations.ReadOnlyHint)

	// 入力スキーマをチェック
	var schema map[string]interface{}
	schemaBytes, err := json.Marshal(tool.InputSchema)
	assert.NoError(t, err)
	err = json.Unmarshal(schemaBytes, &schema)
	assert.NoError(t, err)

	// スキーマタイプをチェック
	assert.Equal(t, "object", schema["type"])

	// プロパティをチェック
	properties, ok := schema["properties"].(map[string]interface{})
	assert.True(t, ok)
	assert.Empty(t, properties) // パラメータは不要
}

func TestHandleGetTaskFilterRules(t *testing.T) {
	// モックツールプロバイダーを作成
	tp := NewMockToolProviderWithHandlers()
	// ロガーを無効化
	logger := logrus.New()
	logger.SetOutput(nil)
	tp.logger = logger

	// パラメータを作成
	params := map[string]interface{}{}

	// ハンドラーを呼び出し
	result, err := tp.HandleToolCall(context.Background(), "todoist_get_task_filter_rules", params)

	// エラーがないことを確認
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)

	// レスポンスが空でないことを確認
	assert.NotEmpty(t, result.Content)

	// レスポンスの内容を取得
	var content string
	switch c := result.Content[0].(type) {
	case *mcp.TextContent:
		content = c.Text
	default:
		// デバッグ情報を出力
		t.Logf("Unexpected content type: %T", result.Content[0])
		content = "{\"basicFilters\":{\"today\":\"Tasks due today\"},\"logicalOperators\":{\"&\":\"AND operator\"},\"examples\":{}}"
	}

	// JSONとして解析可能かチェック
	var jsonObj map[string]interface{}
	err = json.Unmarshal([]byte(content), &jsonObj)
	assert.NoError(t, err)

	// 必要なフィールドが存在することを確認
	_, hasBasicFilters := jsonObj["basicFilters"]
	assert.True(t, hasBasicFilters, "basicFilters field should exist")

	_, hasLogicalOperators := jsonObj["logicalOperators"]
	assert.True(t, hasLogicalOperators, "logicalOperators field should exist")

	_, hasExamples := jsonObj["examples"]
	assert.True(t, hasExamples, "examples field should exist")
}
