package services

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"server/internal/data/db"
	"server/internal/data/storage"
	"time"
)

type DicomService struct {
	minioRepo *storage.MiniRepo
	tasksRepo *db.TasksRepo
}

func NewDicomService(minioRepo *storage.MiniRepo, tasksRepo *db.TasksRepo) *DicomService {
	return &DicomService{
		minioRepo: minioRepo,
		tasksRepo: tasksRepo,
	}
}

func (s *DicomService) GetUrl(id int64) (string, error) {
	task, err := s.tasksRepo.GetTask(id)
	if err != nil {
		return "", err
	}
	objName := fmt.Sprintf("%d", task.ID)
	url, err := s.minioRepo.GetPresignedDownloadURL(context.Background(), objName, time.Hour)
	if err != nil {
		log.Errorf("get presigned download url failed, id [%d], error [%v]", id, err)
		return "", err
	}
	return url, nil
}
