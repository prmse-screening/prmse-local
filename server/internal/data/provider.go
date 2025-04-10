package data

import (
	"context"
	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"server/internal/config"
	"server/internal/data/db"
	"server/internal/data/storage"
	"server/internal/models/entities"
	"time"
)

var ProviderSet = wire.NewSet(NewDatabase, NewMinioClient, db.NewTasksRepo, storage.NewMiniRepo)

func NewDatabase() (*gorm.DB, error) {
	dsn := config.Cfg.Database.Path + "?_journal_mode=WAL&_busy_timeout=5000"

	database, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: &sqlLogger{
			SlowThreshold:             time.Second,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  gormLogger.Warn,
		},
	})
	if err != nil {
		log.Errorf("Failed to connect to database: %s", err.Error())
		return nil, err
	}

	if err = database.AutoMigrate(&entities.Task{}); err != nil {
		log.Errorf("AutoMigrate failed: %s", err.Error())
		return nil, err
	}

	return database, nil
}

func NewMinioClient() (*minio.Client, error) {
	endpoint := config.Cfg.Minio.Endpoint
	accessKeyID := config.Cfg.Minio.AccessKey
	secretAccessKey := config.Cfg.Minio.SecretKey

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Errorf("Failed to create MinIO client: %s", err.Error())
		return nil, err
	}

	bucketName := config.Cfg.Minio.DefaultBucket
	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Errorf("Failed to check if bucket exists: %s", err.Error())
		return nil, err
	}

	if !exists {
		err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Errorf("Failed to create bucket: %s", err.Error())
			return nil, err
		}
		log.Infof("Successfully created bucket %s", bucketName)
	}

	return minioClient, nil
}
