package responses

import (
	"server/internal/commons/enums"
	"time"
)

type CreateTaskResponse struct {
	ID      int64           `json:"id"`
	Series  string          `json:"series"`
	Path    string          `json:"path"`
	Status  enums.TaskState `json:"status"`
	Result  string          `json:"result"`
	Model   string          `json:"model"`
	Order   int64           `json:"order"`
	Updated time.Time       `json:"updated"`
}

type GetUploadUrlResponse struct {
	Url  string            `json:"url"`
	Form map[string]string `json:"form"`
}
