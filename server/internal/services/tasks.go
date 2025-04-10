package services

import "server/internal/data/db"

type TasksService struct {
	tasksRepo *db.TasksRepo
}

func NewTasksService(tasksRepo *db.TasksRepo) *TasksService {
	return &TasksService{tasksRepo: tasksRepo}
}
