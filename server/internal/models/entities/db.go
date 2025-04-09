package entities

import (
	"gorm.io/gorm"
	"time"
)

type Task struct {
	ID      int64     `gorm:"primaryKey;not null;autoIncrement"`
	Series  string    `gorm:"unique;not null"`
	Path    string    `gorm:"not null"`
	Status  string    `gorm:"default:'pending';not null"`
	Result  string    `gorm:"default:NULL"`
	Model   string    `gorm:"not null"`
	Order   int64     `gorm:"not null"`
	Updated time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
}

func (t *Task) BeforeUpdate(_ *gorm.DB) (err error) {
	t.Updated = time.Now()
	return nil
}

func (t *Task) TableName() string {
	return "t_task"
}
