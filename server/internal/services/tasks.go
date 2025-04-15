package services

import (
	"context"
	log "github.com/sirupsen/logrus"
	"server/internal/commons/bizErr"
	"server/internal/config"
	"server/internal/data/db"
	"server/internal/data/storage"
	"server/internal/models/entities"
	"time"
)

type TasksService struct {
	tasksRepo *db.TasksRepo
	minioRepo *storage.MiniRepo
}

func NewTasksService(tasksRepo *db.TasksRepo, miniRepo *storage.MiniRepo) *TasksService {
	return &TasksService{
		tasksRepo: tasksRepo,
		minioRepo: miniRepo,
	}
}

func (s *TasksService) Create(task *entities.Task) error {
	err := s.tasksRepo.Create(task)
	if err != nil {
		log.Errorf("create task failed, series [%s]", task.Series)
		return bizErr.CreateTaskErr
	}
	return nil
}

func (s *TasksService) Update(task *entities.Task) error {
	if err := s.tasksRepo.Update(task); err != nil {
		log.Errorf("update task failed, series [%s]", task.Series)
		return bizErr.UpdateTaskErr
	}
	return nil
}

func (s *TasksService) Delete(task *entities.Task) error {
	err := s.tasksRepo.Delete(task)
	if err != nil {
		log.Errorf("delete task failed, series [%s]", task.Series)
		return bizErr.DeleteTaskErr
	}
	err = s.minioRepo.DeleteFolder(context.Background(), task.Path)
	if err != nil {
		log.Errorf("delete folder failed, path [%s], error [%s]", task.Path, err)
		return bizErr.DeleteTaskErr
	}
	return nil
}

func (s *TasksService) GetUploadUrl(path string) (string, map[string]string, error) {
	url, form, err := s.minioRepo.GetPresignedPostFolderUploadURL(context.Background(), path, time.Hour)
	if err != nil {
		log.Error("get upload url failed, path [%s], error [%s]", path, err)
		return "", nil, err
	}
	return url, form, nil
}

func (s *TasksService) GetListPagination(page, pageSize int, series, status string) ([]*entities.Task, int64, error) {
	tasks, total, err := s.tasksRepo.ListWithPagination(page, pageSize, series, status)
	if err != nil {
		return nil, 0, bizErr.GetTasksErr
	}
	return tasks, total, nil
}

func (s *TasksService) Prioritize(task *entities.Task) error {
	task.Order = task.Order - (time.Hour * 24).Milliseconds()
	if err := s.Update(task); err != nil {
		return bizErr.UpdateTaskErr
	}
	return nil
}

func (s *TasksService) SetWorkerDevice(cpu bool) {
	config.Cfg.Worker.Cpu = cpu
}
