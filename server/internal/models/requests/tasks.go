package requests

import (
	"server/internal/commons/enums"
	"time"
)

type CreateTaskRequest struct {
	Series string `json:"series"`
}

type UpdateTaskRequest struct {
	ID      int64           `json:"id"`
	Series  string          `json:"series"`
	Status  enums.TaskState `json:"status"`
	Result  string          `json:"result"`
	Model   string          `json:"model"`
	Order   int64           `json:"order"`
	Updated time.Time       `json:"updated"`
}

type DeleteTaskRequest struct {
	ID      int64           `json:"id"`
	Series  string          `json:"series"`
	Status  enums.TaskState `json:"status"`
	Result  string          `json:"result"`
	Model   string          `json:"model"`
	Order   int64           `json:"order"`
	Updated time.Time       `json:"updated"`
}
