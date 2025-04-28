package responses

import (
	"gorm.io/datatypes"
	"server/internal/commons/enums"
	"server/internal/models/entities"
	"time"
)

type GetTaskResponse struct {
	ID      int64           `json:"id"`
	Series  string          `json:"series"`
	Status  enums.TaskState `json:"status"`
	Result  datatypes.JSON  `json:"result"`
	Model   string          `json:"model"`
	Order   int64           `json:"order"`
	Updated time.Time       `json:"updated"`
}

type CreateTaskResponse struct {
	ID      int64           `json:"id"`
	Series  string          `json:"series"`
	Status  enums.TaskState `json:"status"`
	Result  datatypes.JSON  `json:"result"`
	Model   string          `json:"model"`
	Order   int64           `json:"order"`
	Updated time.Time       `json:"updated"`
}

type GetUploadUrlResponse struct {
	Url  string            `json:"url"`
	Form map[string]string `json:"form"`
}

type GetTaskLists struct {
	Tasks []*entities.Task `json:"tasks"`
	Total int64            `json:"total"`
}
