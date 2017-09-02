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
	STATUS_CE
	STATUS_RE
	STATUS_TLE
	STATUS_MLE
	STATUS_WA
)
