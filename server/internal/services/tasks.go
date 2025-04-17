package services

import (
	"context"
	log "github.com/sirupsen/logrus"
	"server/internal/commons/bizErr"
	"server/internal/config"
	"server/internal/data/db"
	"server/internal/data/storage"
	"server/internal/models/entities"
	"strconv"
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

	folder := strconv.FormatInt(task.ID, 10)
	err = s.minioRepo.DeleteFolder(context.Background(), folder)
	if err != nil {
		log.Errorf("delete folder failed, path [%s], error [%v]", folder, err)
		return bizErr.DeleteTaskErr
	}
	return nil
}

func (s *TasksService) GetUploadUrl(series string) (string, map[string]string, error) {
	task, err := s.tasksRepo.GetBySeries(series)
	if err != nil {
		log.Errorf("get task failed, series [%s], error [%v]", series, err)
		return "", nil, err
	}
	folder := strconv.FormatInt(task.ID, 10)
	url, form, err := s.minioRepo.GetPresignedPostFolderUploadURL(context.Background(), folder, time.Hour)
	if err != nil {
		log.Errorf("get upload url failed, id [%d], error [%v]", task.ID, err)
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
