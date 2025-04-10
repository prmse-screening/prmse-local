package bizErr

import "errors"

var (
	ParseParamsErr      = errors.New("parse params error")
	RetrieveNextTaskErr = errors.New("retrieve next task failed")
	UpdateTaskErr       = errors.New("update task status failed")
)
