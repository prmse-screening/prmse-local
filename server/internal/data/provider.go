package data

import (
	"github.com/google/wire"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"server/internal/config"
	"server/internal/data/db"
	"server/internal/models/entities"
)

var ProviderSet = wire.NewSet(db.NewTasksRepo, NewDatabase)

func NewDatabase() (*gorm.DB, error) {
	dsn := config.Cfg.Database.Path + "?_journal_mode=WAL&_busy_timeout=5000"
	database, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
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
