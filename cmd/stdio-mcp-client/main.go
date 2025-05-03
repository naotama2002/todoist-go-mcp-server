package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// JSON-RPC リクエスト構造体
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// JSON-RPC レスポンス構造体
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

// JSON-RPC エラー構造体
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// CallToolParams 構造体
type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

func main() {
	// コマンドライン引数からサーバーパスを取得
	serverPath := "cmd/todoist-mcp-server/todoist-mcp-server"
	if len(os.Args) > 1 {
		serverPath = os.Args[1]
	}

	// 環境変数を設定
	env := os.Environ()
	env = append(env, "TODOIST_API_TOKEN="+os.Getenv("TODOIST_API_TOKEN"))

	// サーバーを stdio モードで起動
	cmd := exec.Command(serverPath, "-mode=stdio")
	cmd.Env = env

	// 標準入出力パイプを設定
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Printf("stdin パイプの作成に失敗しました: %v\n", err)
		os.Exit(1)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("stdout パイプの作成に失敗しました: %v\n", err)
		os.Exit(1)
	}

	// コマンドを開始
	if err := cmd.Start(); err != nil {
		fmt.Printf("コマンドの開始に失敗しました: %v\n", err)
		os.Exit(1)
	}

	// 非同期でサーバーからのレスポンスを読み取る
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			// JSON レスポンスを解析
			var response JSONRPCResponse
			if err := json.Unmarshal([]byte(line), &response); err != nil {
				fmt.Printf("サーバーからの出力: %s\n", line)
			} else {
				// 整形された JSON を表示
				var prettyJSON bytes.Buffer
				if err := json.Indent(&prettyJSON, []byte(line), "", "  "); err == nil {
					fmt.Printf("サーバーからのレスポンス:\n%s\n", prettyJSON.String())
				} else {
					fmt.Printf("サーバーからのレスポンス: %s\n", line)
				}
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("stdout の読み取りエラー: %v\n", err)
		}
	}()

	// tools/list リクエストを送信
	listRequest := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/list",
	}

	listRequestJSON, err := json.Marshal(listRequest)
	if err != nil {
		fmt.Printf("リクエストの JSON 変換に失敗しました: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("送信するリクエスト: %s\n", string(listRequestJSON))
	if _, err := stdin.Write(append(listRequestJSON, '\n')); err != nil {
		fmt.Printf("リクエストの送信に失敗しました: %v\n", err)
		os.Exit(1)
	}

	// ユーザー入力を処理
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\n使用可能なコマンド:")
	fmt.Println("list - tools/list リクエストを送信")
	fmt.Println("call <tool_name> [arg1=value1 arg2=value2 ...] - 指定したツールを呼び出す")
	fmt.Println("exit - 終了")
	fmt.Print("> ")

	for scanner.Scan() {
		input := scanner.Text()
		if input == "exit" {
			break
		} else if input == "list" {
			// tools/list リクエストを送信
			listRequest := JSONRPCRequest{
				JSONRPC: "2.0",
				ID:      1,
				Method:  "tools/list",
			}

			listRequestJSON, err := json.Marshal(listRequest)
			if err != nil {
				fmt.Printf("リクエストの JSON 変換に失敗しました: %v\n", err)
				continue
			}

			fmt.Printf("送信するリクエスト: %s\n", string(listRequestJSON))
			if _, err := stdin.Write(append(listRequestJSON, '\n')); err != nil {
				fmt.Printf("リクエストの送信に失敗しました: %v\n", err)
				continue
			}
		} else if strings.HasPrefix(input, "call ") {
			parts := strings.Split(input[5:], " ")
			if len(parts) < 1 {
				fmt.Println("使用法: call <tool_name> [arg1=value1 arg2=value2 ...]")
				continue
			}

			toolName := parts[0]
			arguments := make(map[string]interface{})

			// 引数を解析
			for _, arg := range parts[1:] {
				keyValue := strings.SplitN(arg, "=", 2)
				if len(keyValue) == 2 {
					arguments[keyValue[0]] = keyValue[1]
				}
			}

			// tools/call リクエストを作成
			callToolParams := CallToolParams{
				Name:      toolName,
				Arguments: arguments,
			}

			callRequest := JSONRPCRequest{
				JSONRPC: "2.0",
				ID:      2,
				Method:  "tools/call",
				Params:  callToolParams,
			}

			callRequestJSON, err := json.Marshal(callRequest)
			if err != nil {
				fmt.Printf("リクエストの JSON 変換に失敗しました: %v\n", err)
				continue
			}

			fmt.Printf("送信するリクエスト: %s\n", string(callRequestJSON))
			if _, err := stdin.Write(append(callRequestJSON, '\n')); err != nil {
				fmt.Printf("リクエストの送信に失敗しました: %v\n", err)
				continue
			}
		} else {
			fmt.Println("未知のコマンドです。使用可能なコマンド: list, call, exit")
		}
		fmt.Print("> ")
	}

	// プロセスを終了
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("プロセスの終了に失敗しました: %v\n", err)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Printf("コマンドの終了を待機中にエラーが発生しました: %v\n", err)
	}
}
