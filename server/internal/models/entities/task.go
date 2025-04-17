package entities

import (
	"gorm.io/gorm"
	"server/internal/commons/enums"
	"time"
)

type Task struct {
	ID      int64           `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	Series  string          `gorm:"unique;not null" json:"series"`
	Status  enums.TaskState `gorm:"default:0;not null" json:"status"`
	Result  string          `gorm:"default:null" json:"result"`
	Model   string          `gorm:"not null" json:"model"`
	Order   int64           `gorm:"not null" json:"order"`
	Updated time.Time       `gorm:"not null" json:"updated"`
}

func (t *Task) BeforeUpdate(_ *gorm.DB) (err error) {
	t.Updated = time.Now()
	return nil
}

func (t *Task) BeforeCreate(_ *gorm.DB) (err error) {
	t.Updated = time.Now()
	t.Order = time.Now().UnixMilli()
	return nil
}

func (t *Task) TableName() string {
	return "t_task"
}
