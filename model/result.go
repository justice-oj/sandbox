package model

type Result struct {
	Runtime  int64  `json:"runtime,omitempty"`
	Memory   int64  `json:"memory,omitempty"`
	Status   int32  `json:"status"`
	Error    string `json:"error,omitempty"`
	Input    string `json:"input,omitempty"`
	Output   string `json:"output,omitempty"`
	Expected string `json:"expected,omitempty"`
}

const (
	StatusAc = iota
	_
	StatusRe
	StatusTle
	_
	StatusWa
)

func (r *Result) GetAcceptedTaskResult(runtime, memory int64) *Result {
	r.Status = StatusAc
	r.Runtime = runtime
	r.Memory = memory
	return r
}

func (r *Result) GetRuntimeErrorTaskResult() *Result {
	r.Status = StatusRe
	r.Error = "Runtime Error"
	return r
}

func (r *Result) GetTimeLimitExceededErrorTaskResult() *Result {
	r.Status = StatusTle
	r.Error = "Runtime Error"
	return r
}

func (r *Result) GetWrongAnswerTaskResult(input, output, expected string) *Result {
	r.Status = StatusWa
	r.Input = input
	r.Output = output
	r.Expected = expected
	return r
}
