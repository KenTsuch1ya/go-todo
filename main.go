package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Done bool   `json:done`
}

const taskFile = "tasks.json"

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("使用方法: todo <command> [arguments]")
		fmt.Println("利用可能なコマンド: add <task>, list")
		return
	}

	command := args[1]

	switch command {
	case "add":
		if len(args) < 3 {
			fmt.Println("タスクの内容を入力してください")
			return
		}
		task := args[2]
		addTask(task)

	case "list":
		listTasks()

	case "done":
		if len(args) < 3 {
			fmt.Println("完了するタスクのIDを入力してください")
			return
		}
		id := args[2]
		markTaskDone(id)

	case "remove":
		if len(args) < 3 {
			fmt.Println("削除するタスクのIDを入力してください")
			return
		}
		id := args[2]
		removeTask(id)

	default:
		fmt.Println("不明なコマンドです。利用可能なコマンド: add, list")
	}
}

func addTask(task string) {
	tasks := readTasks()
	newTask := Task{ID: len(tasks) + 1, Name: task}
	tasks = append(tasks, newTask)
	saveTasks(tasks)
	fmt.Printf("タスク追加：%s\n", task)
}

func listTasks() {
	tasks := readTasks()
	if len(tasks) == 0 {
		fmt.Println("タスクはありません。")
		return
	}
	fmt.Println("タスク一覧：")
	for _, task := range tasks {
		status := "[ ]"
		if task.Done {
			status = "[✔]"
		}
		fmt.Printf("%s [%d] %s\n", status, task.ID, task.Name)
	}
}

func markTaskDone(idStr string) {
	tasks := readTasks()
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("エラー: JSONの読み込みに失敗しました。", err)
		return
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = true
			saveTasks(tasks)
			fmt.Printf("タスク完了: %s\n", tasks[i].Name)
			return
		}
	}
	fmt.Println("エラー: 指定したIDのタスクが見つかりません")
}

func removeTask(idStr string) {
	tasks := readTasks()
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("エラー: IDは数字で入力してください")
		return
	}

	newTasks := []Task{}
	for _, task := range tasks {
		if task.ID != id {
			newTasks = append(newTasks, task)
		}
	}

	if len(newTasks) == len(tasks) {
		fmt.Println("エラー: 指定したIDのタスクが見つかりません")
		return
	}

	for i := range newTasks {
		newTasks[i].ID = i + 1
	}

	saveTasks(newTasks)
	fmt.Println(("タスクを削除しました"))
}

func saveTasks(tasks []Task) {
	file, err := os.Create(taskFile)
	if err != nil {
		fmt.Println("エラー：タスクを保存できませんでした。", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(tasks)
	if err != nil {
		fmt.Println("エラー: JSONの書き込みに失敗しました。", err)
	}
}

func readTasks() []Task {
	if _, err := os.Stat(taskFile); os.IsNotExist(err) {
		return []Task{}
	}

	file, err := os.Open(taskFile)
	if err != nil {
		fmt.Println("エラー: タスクを読み込めませんでした。", err)
		return []Task{}
	}
	defer file.Close()

	var tasks []Task
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil {
		fmt.Println("エラー: JSONの読み込みに失敗しました。", err)
		return []Task{}
	}

	return tasks
}
