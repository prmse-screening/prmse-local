package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"server/internal/models/entities"
)

var allowedSortKeys = map[string]bool{"order": true, "updated": true}

type TasksRepo struct {
	db *gorm.DB
}

func NewTasksRepo(db *gorm.DB) *TasksRepo {
	return &TasksRepo{db: db}
}

func (r *TasksRepo) Create(task *entities.Task) error {
	return r.db.Create(task).Error
}

func (r *TasksRepo) GetTask(id int64) (*entities.Task, error) {
	var task entities.Task
	err := r.db.Where("id = ?", id).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
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

func (r *TasksRepo) Delete(task *entities.Task) error {
	return r.db.Delete(task).Error
}

func (r *TasksRepo) ListWithPagination(page, pageSize int, status, series, sortKey, sortOrder string) ([]*entities.Task, int64, error) {
	var tasks []*entities.Task
	var total int64

	query := r.db.Model(&entities.Task{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if series != "" {
		query = query.Where("series LIKE ?", series+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if allowedSortKeys[sortKey] {
		orderClause := clause.OrderByColumn{Column: clause.Column{Name: sortKey}, Desc: sortOrder == "desc"}
		query = query.Order(orderClause)
	}

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	if err := query.Limit(pageSize).Offset(offset).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *TasksRepo) NextTask() (*entities.Task, error) {
	var task entities.Task
	err := r.db.Where("status = ?", 1).
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
