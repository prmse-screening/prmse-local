package services

import (
	"context"
	"encoding/csv"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"server/internal/commons/bizErr"
	"server/internal/commons/enums"
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

func (s *TasksService) GetTask(id int64) (*entities.Task, error) {
	task, err := s.tasksRepo.GetTask(id)
	if err != nil {
		log.Errorf("get task failed, id [%d], error [%v]", id, err)
		return nil, bizErr.GetTaskErr
	}
	return task, nil
}

func (s *TasksService) ExportTasks(series, status string) (string, error) {
	rows, err := s.tasksRepo.GetCursor(series, status)
	if err != nil {
		log.Errorf("get next tasks failed, error [%v]", err)
		return "", bizErr.GetTasksErr
	}
	defer func() {
		_ = rows.Close()
	}()

	filename := fmt.Sprintf("tasks_export_%s.csv", time.Now().String())
	filePath := path.Join("exports", filename)

	if err = os.MkdirAll("exports", os.ModePerm); err != nil {
		log.Errorf("failed to create exports directory: %v", err)
		return "", err
	}

	file, err := os.Create(filePath)
	if err != nil {
		log.Errorf("failed to create csv file: %v", err)
		return "", err
	}
	writer := csv.NewWriter(file)
	defer func() {
		_ = file.Close()
	}()
	defer writer.Flush()

	err = writer.Write([]string{"ID", "Series", "Status", "Result", "Model", "Order", "Updated"})
	if err != nil {
		log.Errorf("failed to write header to csv: %v", err)
		return "", err
	}

	for rows.Next() {
		var task entities.Task
		if err = rows.Scan(
			&task.ID,
			&task.Series,
			&task.Status,
			&task.Result,
			&task.Model,
			&task.Order,
			&task.Updated,
		); err != nil {
			log.Errorf("scan task failed: %v", err)
			continue
		}
		err = writer.Write([]string{
			fmt.Sprintf("%d", task.ID),
			task.Series,
			task.Status.String(),
			task.Result,
			task.Model,
			fmt.Sprintf("%d", task.Order),
			task.Updated.Format(time.RFC3339),
		})
		if err != nil {
			log.Errorf("failed to write task to csv: %v", err)
			return "", err
		}
	}
	return filePath, nil
}

func (s *TasksService) Create(task *entities.Task) (*entities.Task, error) {
	err := s.tasksRepo.Create(task)
	if err != nil {
		log.Errorf("create task failed, series [%s]: %v", task.Series, err)

		existing, findErr := s.tasksRepo.GetBySeries(task.Series)
		if findErr == nil && existing.Status == enums.Preparing {
			return existing, nil
		}

		return nil, bizErr.CreateTaskErr
	}

	return task, nil
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

func (s *TasksService) GetListPagination(page, pageSize int, series, status, sortKey, sortOrder string) ([]*entities.Task, int64, error) {
	tasks, total, err := s.tasksRepo.ListWithPagination(page, pageSize, series, status, sortKey, sortOrder)
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
