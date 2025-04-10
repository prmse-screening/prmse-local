package db

import (
	"gorm.io/gorm"
	"server/internal/models/entities"
)

type TasksRepo struct {
	db *gorm.DB
}

func NewTasksRepo(db *gorm.DB) *TasksRepo {
	return &TasksRepo{db: db}
}

func (r *TasksRepo) Create(task *entities.Task) error {
	return r.db.Create(task).Error
}

func (r *TasksRepo) GetBySeries(series string) (*entities.Task, error) {
	var task entities.Task
	err := r.db.Where("series = ?", series).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TasksRepo) Update(task *entities.Task) error {
	return r.db.Save(task).Error
}

func (r *TasksRepo) Delete(id int64) error {
	return r.db.Delete(&entities.Task{}, id).Error
}

func (r *TasksRepo) ListWithPagination(page, pageSize int, status, series string) ([]*entities.Task, int64, error) {
	var tasks []*entities.Task
	var total int64

	query := r.db

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if series != "" {
		query = query.Where("series LIKE ?", series+"%")
	}

	err := query.Model(&entities.Task{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("\"order\"").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&tasks).Error

	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *TasksRepo) NextTask() (*entities.Task, error) {
	var task entities.Task
	err := r.db.Where("status = ?", "pending").
		Order("\"order\"").Take(&task).Error

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *TasksRepo) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&entities.Task{}).Where("status = ?", status).Count(&count).Error
	return count, err
}
