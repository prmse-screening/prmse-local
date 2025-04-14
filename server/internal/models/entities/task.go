package entities

import (
	"gorm.io/gorm"
	"server/internal/commons/enums"
	"time"
)

type Task struct {
	ID      int64           `gorm:"primaryKey;not null;autoIncrement"`
	Series  string          `gorm:"unique;not null"`
	Path    string          `gorm:"not null"`
	Status  enums.TaskState `gorm:"default:0;not null"`
	Result  string          `gorm:"default:null"`
	Model   string          `gorm:"not null"`
	Order   int64           `gorm:"not null"`
	Updated time.Time       `gorm:"not null"`
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
