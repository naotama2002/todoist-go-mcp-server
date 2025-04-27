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
	serverURL := "http://localhost:8081"

	// todoist_get_tasks ツールを呼び出すリクエストを作成
	callToolParams := CallToolParams{
		Name:      "todoist_get_tasks",
		Arguments: map[string]interface{}{},
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
	defer resp.Body.Close()

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
