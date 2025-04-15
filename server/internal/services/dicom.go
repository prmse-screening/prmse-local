package services

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"server/internal/data/storage"
	"time"
)

type DicomService struct {
	minioRepo *storage.MiniRepo
}

func NewDicomService(minioRepo *storage.MiniRepo) *DicomService {
	return &DicomService{
		minioRepo: minioRepo,
	}
}

func (s *DicomService) Redirect(series string, file string) (string, error) {
	objName := fmt.Sprintf("%s/%s", series, file)
	url, err := s.minioRepo.GetPresignedDownloadURL(context.Background(), objName, time.Hour)
	if err != nil {
		log.Errorf("get presigned download url failed, series [%s], file [%s], error [%s]", series, file, err)
		return "", err
	}
	return url, nil
}
