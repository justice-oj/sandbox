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
	STATUS_AC  = iota
	_
	STATUS_RE
	STATUS_TLE
	_
	STATUS_WA
)

func (r *Result) GetAcceptedTaskResult(runtime, memory int64) *Result {
	r.Status = STATUS_AC
	r.Runtime = runtime
	r.Memory = memory
	return r
}

func (r *Result) GetRuntimeErrorTaskResult() *Result {
	r.Status = STATUS_RE
	r.Error = "Runtime Error"
	return r
}

func (r *Result) GetTimeLimitExceededErrorTaskResult() *Result {
	r.Status = STATUS_TLE
	r.Error = "Runtime Error"
	return r
}

func (r *Result) GetWrongAnswerTaskResult(input, output, expected string) *Result {
	r.Status = STATUS_WA
	r.Input = input
	r.Output = output
	r.Expected = expected
	return r
}
