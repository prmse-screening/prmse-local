package schedule

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

type TasksCleaner struct{}

func NewTaskCleaner() *TasksCleaner {
	return &TasksCleaner{}
}

func (c *TasksCleaner) StartCSVFileCleaner() {
	ticker := time.NewTicker(time.Minute)
	go func() {
		for range ticker.C {
			c.cleanupOldCSVFiles("exports", time.Minute)
		}
	}()
}

func (c *TasksCleaner) cleanupOldCSVFiles(dir string, maxAge time.Duration) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Errorf("readDir failed: %v", err)
		return
	}

	now := time.Now()
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := path.Join(dir, file.Name())
		info, err := os.Stat(filePath)
		if err != nil {
			log.Warnf("stat file failed: %v", err)
			continue
		}
		if now.Sub(info.ModTime()) > maxAge {
			if err := os.Remove(filePath); err != nil {
				log.Warnf("failed to remove file %s: %v", filePath, err)
			} else {
				log.Infof("removed old CSV file: %s", filePath)
			}
		}
	}
}
