package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/internal/models/entities"
	"server/internal/models/responses"
	"server/internal/services"
	"server/internal/utils"
)

type TasksHandler struct {
	tasksService *services.TasksService
}

func NewTasksHandler(tasksService *services.TasksService) *TasksHandler {
	return &TasksHandler{
		tasksService: tasksService,
	}
}

func (h *TasksHandler) CreateTask(c *gin.Context) {
	var task entities.Task

	if err := utils.Bind(c, &task); err != nil {
		return
	}

	var resp responses.BaseResponse
	if err := h.tasksService.Create(); err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp.SuccessResponse(c, http.StatusCreated, task)
}
