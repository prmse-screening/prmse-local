package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/internal/commons/bizErr"
	"server/internal/models/requests"
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
	var req requests.CreateTaskRequest

	if err := utils.Bind(c, &req); err != nil {
		return
	}

	var resp responses.BaseResponse
	data, err := h.tasksService.Create(&req)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	resp.SuccessResponse(c, http.StatusOK, data)
}

func (h *TasksHandler) UpdateTask(c *gin.Context) {
	var req requests.UpdateTaskRequest

	if err := utils.Bind(c, &req); err != nil {
		return
	}

	var resp responses.BaseResponse
	if err := h.tasksService.Update(&req); err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	resp.SuccessResponse(c, http.StatusOK, nil)
}

func (h *TasksHandler) DeleteTask(c *gin.Context) {
	var req requests.DeleteTaskRequest

	if err := utils.Bind(c, &req); err != nil {
		return
	}

	var resp responses.BaseResponse
	if err := h.tasksService.Delete(&req); err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	resp.SuccessResponse(c, http.StatusOK, nil)
}

func (h *TasksHandler) GetUploadUrl(c *gin.Context) {
	path := c.Query("path")
	var resp responses.BaseResponse
	if path == "" {
		resp.ErrorResponse(c, http.StatusBadRequest, bizErr.ParseParamsErr)
		return
	}
	data, err := h.tasksService.GetUploadUrl(path)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	resp.SuccessResponse(c, http.StatusOK, data)
}
