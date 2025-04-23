package schedule

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"server/internal/commons/bizErr"
	"server/internal/commons/enums"
	"server/internal/config"
	"server/internal/data/db"
	"server/internal/data/storage"
	"server/internal/rpc"
	"strconv"
	"sync/atomic"
	"time"
)

type TasksScheduler struct {
	tasksRepo       *db.TasksRepo
	minioRepo       *storage.MiniRepo
	workers         *[]rpc.WorkerClient
	semaphore       chan int
	ctx             context.Context
	cancelFunc      context.CancelFunc
	nextWorkerIndex uint64
}

func NewTasksScheduler(tasksRepo *db.TasksRepo, workers *[]rpc.WorkerClient, minioRepo *storage.MiniRepo) *TasksScheduler {
	nums := len(*workers)
	semaphore := make(chan int, nums)
	for i := 0; i < nums; i++ {
		semaphore <- i
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &TasksScheduler{
		tasksRepo:       tasksRepo,
		minioRepo:       minioRepo,
		workers:         workers,
		semaphore:       semaphore,
		ctx:             ctx,
		cancelFunc:      cancelFunc,
		nextWorkerIndex: 0,
	}
}

func (s *TasksScheduler) processNextTask(id int) {
	defer func() {
		s.semaphore <- id
	}()
	if err := s.handleNextTask(id); err != nil {
		if !errors.Is(err, bizErr.RetrieveNextTaskErr) {
			log.Errorf("%v", err)
		}
		time.Sleep(6 * time.Second)
	}
}

func (s *TasksScheduler) handleNextTask(id int) error {
	task, err := s.tasksRepo.NextTask()
	if err != nil {
		return bizErr.RetrieveNextTaskErr
	}

	task.Status = enums.Processing
	if err = s.tasksRepo.Update(task); err != nil {
		return bizErr.UpdateTaskErr
	}

	worker := (*s.workers)[id]
	//reqCtx, cancel := context.WithTimeout(s.ctx, time.Hour)
	//defer cancel()

	url, err := s.minioRepo.GetPresignedDownloadURL(s.ctx, strconv.FormatInt(task.ID, 10), time.Hour)
	if err != nil {
		return bizErr.GetDownloadUrlsErr
	}
	resp, err := worker.Infer(s.ctx, &rpc.InferenceRequest{
		Model:  task.Model,
		Path:   url,
		Series: task.Series,
		Cpu:    config.Cfg.Worker.Cpu,
	})

	if err != nil {
		task.Status = enums.Failed
		_ = s.tasksRepo.Update(task)
		return err
	}

	task.Status = enums.Completed
	task.Result = resp.Result

	if err = s.tasksRepo.Update(task); err != nil {
		return bizErr.UpdateTaskErr
	}

	return nil
}

func (s *TasksScheduler) selectWorker() rpc.WorkerClient {
	workers := *s.workers
	workerIndex := atomic.AddUint64(&s.nextWorkerIndex, 1) % uint64(len(workers))
	return workers[workerIndex]
}

func (s *TasksScheduler) Start() {
	s.ctx, s.cancelFunc = context.WithCancel(context.Background())

	for {
		select {
		case <-s.ctx.Done():
			return
		case i, ok := <-s.semaphore:
			if !ok {
				return
			}
			go s.processNextTask(i)
		}
	}
}

func (s *TasksScheduler) Stop() {
	s.cancelFunc()
}
