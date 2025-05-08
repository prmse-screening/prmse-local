package data

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
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
	var database *gorm.DB
	var err error

	switch config.Cfg.Database.Source {
	case "MySQL":
		{
			database, err = newMySQL()
		}
	case "SQLite":
		{
			database, err = newSQLite()
		}
	default:
		{
			database, err = newMySQL()
		}
	}
	if err != nil {
		log.Errorf("Failed to connect to %s database: %v", config.Cfg.Database.Source, err)
		return nil, err
	}

	if err = database.AutoMigrate(&entities.Task{}); err != nil {
		log.Errorf("Failed to migrate database: %v", err)
		return nil, err
	}

	return database, nil
}

func newSQLite() (*gorm.DB, error) {
	dsn := config.Cfg.Database.SQLite.Path + "?_journal_mode=WAL&_busy_timeout=5000"

	database, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: &sqlLogger{
			SlowThreshold:             time.Second,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  gormLogger.Warn,
		},
	})
	if err != nil {
		return nil, err
	}
	return database, nil
}

func newMySQL() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Cfg.Database.MySQL.Username,
		config.Cfg.Database.MySQL.Password,
		config.Cfg.Database.MySQL.Host,
		config.Cfg.Database.MySQL.Port,
		config.Cfg.Database.MySQL.Name,
	)
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger: &sqlLogger{
			SlowThreshold:             time.Second,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  gormLogger.Warn,
		},
	})

	if err != nil {
		return nil, err
	}

	if sqlDB, err := database.DB(); err == nil && sqlDB != nil {
		sqlDB.SetMaxIdleConns(config.Cfg.Database.MySQL.MaxIdleConns)
		sqlDB.SetMaxOpenConns(config.Cfg.Database.MySQL.SetMaxOpenConns)
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
		log.Errorf("Failed to create MinIO client: %v", err)
		return nil, err
	}

	bucketName := config.Cfg.Minio.DefaultBucket
	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Errorf("Failed to check if bucket exists: %v", err)
		return nil, err
	}

	if !exists {
		err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Errorf("Failed to create bucket: %v", err)
			return nil, err
		}
		log.Infof("Successfully created bucket %s", bucketName)
	}

	return minioClient, nil
}
