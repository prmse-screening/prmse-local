package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/internal/data/db"
	"server/internal/utils"
)

type TasksHandler struct {
	TasksRepo *db.TasksRepo
}

func NewTasksHandler(tasksRepo *db.TasksRepo) *TasksHandler {
	return &TasksHandler{
		TasksRepo: tasksRepo,
	}
}

func (h *TasksHandler) CreateTask(c *gin.Context) {
	var task db.Task

	if err := utils.Bind(c, &task); err != nil {
		return
	}

	if err := h.TasksRepo.Create(&task); err != nil {
		responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	responses.SuccessResponse(c, http.StatusCreated, task)
}
