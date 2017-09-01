package models

type TaskResult struct {
	Runtime  int64
	Memory   int64
	Status   int32
	Error    string
	Input    string
	Output   string
	Expected string
}
