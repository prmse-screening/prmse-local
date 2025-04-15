package bizErr

import "errors"

var (
	ParseParamsErr = errors.New("parse params error")
)

// task error
var (
	RetrieveNextTaskErr = errors.New("retrieve next task failed")
	CreateTaskErr       = errors.New("create task failed")
	UpdateTaskErr       = errors.New("update task failed")
	DeleteTaskErr       = errors.New("delete task failed")
	GetPostUrlErr       = errors.New("get post url failed")
	GetDownloadUrlsErr  = errors.New("get download urls failed")
	GetTasksErr         = errors.New("get tasks failed")
)
