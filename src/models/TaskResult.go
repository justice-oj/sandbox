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

const (
	STATUS_AC  = iota
	_
	STATUS_RE
	STATUS_TLE
	STATUS_MLE
	STATUS_WA
)

func GetAccepptedTaskResult(runtime, memory int64) TaskResult {
	return TaskResult{
		Status:  STATUS_AC,
		Runtime: runtime,
		Memory:  memory,
	}
}

func GetRuntimeErrorTaskResult() TaskResult {
	return TaskResult{
		Status: STATUS_RE,
		Error:  "Runtime Error",
	}
}

func GetTimeLimitExceededErrorTaskResult() TaskResult {
	return TaskResult{
		Status: STATUS_TLE,
		Error:  "Runtime Error",
	}
}

func GetMemoryLimitExceededErrorTaskResult() TaskResult {
	return TaskResult{
		Status: STATUS_MLE,
		Error:  "Runtime Error",
	}
}

func GetWrongAnswerTaskResult(input, output, expected string) TaskResult {
	return TaskResult{
		Status:   STATUS_WA,
		Input:    input,
		Output:   output,
		Expected: expected,
	}
}
