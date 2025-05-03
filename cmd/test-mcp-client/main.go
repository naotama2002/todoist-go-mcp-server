package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// JSON-RPC リクエスト構造体
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

// CallToolParams 構造体
type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

func main() {
	// MCP サーバーの URL
	serverURL := "http://localhost:8082"

	// コマンドライン引数からツール名を取得
	toolName := "todoist_get_tasks"
	if len(os.Args) > 1 {
		toolName = os.Args[1]
	}

	// ツール引数を初期化
	arguments := map[string]interface{}{}

	// ツール名に応じて引数を設定
	switch toolName {
	case "todoist_get_tasks":
		// 引数なし（デフォルト）
	case "todoist_get_task":
		if len(os.Args) < 3 {
			fmt.Println("使用法: go run main.go todoist_get_task <task_id>")
			os.Exit(1)
		}
		arguments["id"] = os.Args[2]
	case "todoist_create_task":
		if len(os.Args) < 3 {
			fmt.Println("使用法: go run main.go todoist_create_task <content>")
			os.Exit(1)
		}
		arguments["content"] = os.Args[2]
		if len(os.Args) > 3 {
			arguments["description"] = os.Args[3]
		}
	case "todoist_update_task":
		if len(os.Args) < 4 {
			fmt.Println("使用法: go run main.go todoist_update_task <task_id> <content>")
			os.Exit(1)
		}
		arguments["id"] = os.Args[2]
		arguments["content"] = os.Args[3]
	case "todoist_close_task":
		if len(os.Args) < 3 {
			fmt.Println("使用法: go run main.go todoist_close_task <task_id>")
			os.Exit(1)
		}
		arguments["id"] = os.Args[2]
	case "todoist_delete_task":
		if len(os.Args) < 3 {
			fmt.Println("使用法: go run main.go todoist_delete_task <task_id>")
			os.Exit(1)
		}
		arguments["id"] = os.Args[2]
	case "todoist_get_projects":
		// 引数なし（デフォルト）
	case "todoist_get_project":
		if len(os.Args) < 3 {
			fmt.Println("使用法: go run main.go todoist_get_project <project_id>")
			os.Exit(1)
		}
		arguments["id"] = os.Args[2]
	default:
		fmt.Printf("未知のツール名: %s\n", toolName)
		os.Exit(1)
	}

	// todoist_get_tasks ツールを呼び出すリクエストを作成
	callToolParams := CallToolParams{
		Name:      toolName,
		Arguments: arguments,
	}

	// JSON-RPC リクエストを作成
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/call",
		Params:  callToolParams,
	}

	// リクエストを JSON に変換
	requestJSON, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("リクエストの JSON 変換に失敗しました: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("送信するリクエスト: %s\n", string(requestJSON))

	// HTTP リクエストを送信
	resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(requestJSON))
	if err != nil {
		fmt.Printf("HTTP リクエストの送信に失敗しました: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("レスポンスボディのクローズに失敗しました: %v\n", err)
		}
	}()

	// レスポンスを読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("レスポンスの読み取りに失敗しました: %v\n", err)
		os.Exit(1)
	}

	// レスポンスを表示
	fmt.Printf("ステータスコード: %d\n", resp.StatusCode)
	fmt.Printf("レスポンス: %s\n", string(body))

	// レスポンスを整形して表示
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "  "); err == nil {
		fmt.Printf("整形されたレスポンス:\n%s\n", prettyJSON.String())
	}
}
