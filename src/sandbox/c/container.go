package main

import (
	"encoding/json"
	"../../models"
	"os"
)

func main() {
	task_result := &models.TaskResult{}

	result, _ := json.Marshal(task_result)
	os.Stdout.Write(result)
}
