package main

import (
	"fmt"
	"os"

	"github.com/naotama2002/todoist-go-mcp-server/pkg/log"
	"github.com/naotama2002/todoist-go-mcp-server/pkg/todoist"
)

func main() {
	// トークンを環境変数から取得
	token := os.Getenv("TODOIST_API_TOKEN")
	if token == "" {
		fmt.Println("TODOIST_API_TOKEN 環境変数が設定されていません")
		os.Exit(1)
	}

	// ロガーを作成
	logger := log.NewLogger()

	// クライアントを作成
	client := todoist.NewClient(token, todoist.WithLogger(logger))

	// タスクを取得
	tasks, err := client.GetTasks("", "")
	if err != nil {
		logger.WithError(err).Error("タスクの取得に失敗しました")
		os.Exit(1)
	}

	// 結果を表示
	fmt.Printf("取得したタスク数: %d\n", len(tasks))
	for i, task := range tasks {
		if i >= 5 {
			fmt.Println("...")
			break
		}
		fmt.Printf("- ID: %s, Content: %s\n", task.ID, task.Content)
	}
}
