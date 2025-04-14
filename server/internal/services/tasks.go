package services

import (
	"context"
	"github.com/jinzhu/copier"
	log "github.com/sirupsen/logrus"
	"server/internal/commons/bizErr"
	"server/internal/data/db"
	"server/internal/data/storage"
	"server/internal/models/entities"
	"server/internal/models/requests"
	"server/internal/models/responses"
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

func (s *TasksService) Create(req *requests.CreateTaskRequest) (*responses.CreateTaskResponse, error) {
	task := entities.Task{
		Series: req.Series,
		Model:  "Sybil",
		Path:   req.Series,
	}
	err := s.tasksRepo.Create(&task)
	if err != nil {
		log.Errorf("create task failed, series [%s]", req.Series)
		return nil, bizErr.CreateTaskErr
	}

	resp := responses.CreateTaskResponse{}
	_ = copier.Copy(&resp, task)

	return &resp, nil
}

func (s *TasksService) Update(req *requests.UpdateTaskRequest) error {
	var task entities.Task
	_ = copier.Copy(&task, req)
	err := s.tasksRepo.Update(&task)
	if err != nil {
		log.Errorf("update task failed, series [%s]", req.Series)
		return bizErr.UpdateTaskErr
	}
	return nil
}

func (s *TasksService) Delete(req *requests.DeleteTaskRequest) error {
	var task entities.Task
	_ = copier.Copy(&task, req)
	err := s.tasksRepo.Delete(&task)
	if err != nil {
		log.Errorf("delete task failed, series [%s]", req.Series)
		return bizErr.DeleteTaskErr
	}
	err = s.minioRepo.DeleteFolder(context.Background(), req.Path)
	if err != nil {
		log.Errorf("delete folder failed, path [%s], error [%s]", req.Path, err)
		return bizErr.DeleteTaskErr
	}
	return nil
}

func (s *TasksService) GetUploadUrl(path string) (*responses.GetUploadUrlResponse, error) {
	url, form, err := s.minioRepo.GetPresignedPostFolderUploadURL(context.Background(), path, time.Hour)
	if err != nil {
		log.Error("get upload url failed, path [%s], error [%s]", path, err)
		return nil, bizErr.GetPostUrlErr
	}
	return &responses.GetUploadUrlResponse{
		Url:  url,
		Form: form,
	}, nil
}

func (s *TasksService) ListP(req *requests.UpdateTaskRequest) error {
	var task entities.Task
	_ = copier.Copy(&task, req)
	err := s.tasksRepo.Update(&task)
	if err != nil {
		return bizErr.UpdateTaskErr
	}
	return nil
}
