package schedule

import (
	"context"
	log "github.com/sirupsen/logrus"
	"server/internal/constants/bizErr"
	"server/internal/data/db"
	"server/internal/data/storage"
	"server/internal/rpc"
	"sync/atomic"
	"time"
)

type TasksScheduler struct {
	tasksRepo       *db.TasksRepo
	workers         *[]rpc.WorkerClient
	minioRepo       *storage.MiniRepo
	workerSem       chan struct{}
	ctx             context.Context
	cancelFunc      context.CancelFunc
	nextWorkerIndex uint64
}

func NewTasksScheduler(tasksRepo *db.TasksRepo, workers *[]rpc.WorkerClient, minioRepo *storage.MiniRepo) *TasksScheduler {
	return &TasksScheduler{
		tasksRepo:       tasksRepo,
		workers:         workers,
		minioRepo:       minioRepo,
		workerSem:       make(chan struct{}, len(*workers)),
		ctx:             nil,
		cancelFunc:      nil,
		nextWorkerIndex: 0,
	}
}

func (s *TasksScheduler) processNextTask() {
	defer func() {
		s.workerSem <- struct{}{}
	}()

	if err := s.handleNextTask(); err != nil {
		log.Error(err)
		time.Sleep(5 * time.Second)
	}
}

func (s *TasksScheduler) handleNextTask() error {
	task, err := s.tasksRepo.NextTask()
	if err != nil {
		return bizErr.RetrieveNextTaskErr
	}

	task.Status = "Processing"
	if err = s.tasksRepo.Update(task); err != nil {
		return bizErr.UpdateTaskErr
	}

	worker := s.selectWorker()
	reqCtx, cancel := context.WithTimeout(s.ctx, 10*time.Minute)
	defer cancel()

	resp, err := worker.Inference(reqCtx, &rpc.InferenceRequest{
		Model:  task.Model,
		Path:   task.Path,
		Series: task.Series,
	})

	if err != nil {
		task.Status = "failed"
	} else {
		task.Status = "completed"
		task.Result = resp.Result
	}

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

	for i := 0; i < cap(s.workerSem); i++ {
		s.workerSem <- struct{}{}
	}

	for {
		select {
		case <-s.ctx.Done():
			return
		case _, ok := <-s.workerSem:
			if !ok {
				return
			}
			go s.processNextTask()
		}
	}
}

func (s *TasksScheduler) Stop() {
	s.cancelFunc()
}
