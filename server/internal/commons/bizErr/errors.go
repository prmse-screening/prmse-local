package bizErr

import "errors"

var (
	ParseParamsErr = errors.New("parse params error")
)

// task error
var (
	RetrieveNextTaskErr = errors.New("RetrieveNextTaskErr")
	InferenceErr        = errors.New("inference failed")
	GetTaskErr          = errors.New("get task failed")
	CreateTaskErr       = errors.New("create task failed, task may already exist")
	UpdateTaskErr       = errors.New("update task failed")
	DeleteTaskErr       = errors.New("delete task failed")
	GetPostUrlErr       = errors.New("get post url failed")
	GetDownloadUrlsErr  = errors.New("get download urls failed")
	GetTasksErr         = errors.New("get tasks failed")
)
