package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"server/internal/commons/bizErr"
	"server/internal/models/entities"
	"server/internal/models/requests"
	"server/internal/models/responses"
	"server/internal/services"
	"server/internal/utils"
	"strconv"
	"strings"
	"time"
)

type TasksHandler struct {
	tasksService *services.TasksService
}

func NewTasksHandler(tasksService *services.TasksService) *TasksHandler {
	return &TasksHandler{
		tasksService: tasksService,
	}
}

func (h *TasksHandler) GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	var resp responses.BaseResponse
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, bizErr.ParseParamsErr.Error())
		return
	}

	task, err := h.tasksService.GetTask(id)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	data := responses.GetTaskResponse{}
	_ = copier.Copy(&data, &task)
	resp.SuccessResponse(c, http.StatusOK, data)
}

func (h *TasksHandler) ExportTasks(c *gin.Context) {
	status := strings.TrimSpace(c.Query("status"))
	series := strings.TrimSpace(c.Query("series"))
	filePath, err := h.tasksService.ExportTasks(series, status)
	var resp responses.BaseResponse
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.FileAttachment(filePath, fmt.Sprintf("tasks_%s.csv", time.Now().String()))
}

func (h *TasksHandler) CreateTask(c *gin.Context) {
	var req requests.CreateTaskRequest
	if err := utils.Bind(c, &req); err != nil {
		return
	}
	task := entities.Task{
		Series: req.Series,
		Model:  "Sybil",
	}

	var resp responses.BaseResponse
	newTask, err := h.tasksService.Create(&task)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	data := responses.CreateTaskResponse{}
	_ = copier.Copy(&data, &newTask)
	resp.SuccessResponse(c, http.StatusOK, data)
}

func (h *TasksHandler) UpdateTask(c *gin.Context) {
	var req requests.UpdateTaskRequest
	if err := utils.Bind(c, &req); err != nil {
		return
	}
	var task entities.Task
	_ = copier.Copy(&task, req)

	var resp responses.BaseResponse
	if err := h.tasksService.Update(&task); err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp.SuccessResponse(c, http.StatusOK, "success")
}

func (h *TasksHandler) PrioritizeTask(c *gin.Context) {
	var req requests.UpdateTaskRequest
	if err := utils.Bind(c, &req); err != nil {
		return
	}
	var task entities.Task
	_ = copier.Copy(&task, req)

	var resp responses.BaseResponse
	if err := h.tasksService.Prioritize(&task); err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp.SuccessResponse(c, http.StatusOK, "success")
}

func (h *TasksHandler) GetListPagination(c *gin.Context) {
	pageParam := c.Query("page")
	pageSizeParam := c.Query("pageSize")
	status := strings.TrimSpace(c.Query("status"))
	series := strings.TrimSpace(c.Query("series"))
	sortKey := strings.TrimSpace(c.Query("sortKey"))
	sortOrder := strings.ToLower(strings.TrimSpace(c.Query("sortOrder")))

	var resp responses.BaseResponse
	if pageParam == "" || pageSizeParam == "" {
		resp.ErrorResponse(c, http.StatusBadRequest, bizErr.ParseParamsErr.Error())
		return
	}
	page, _ := strconv.Atoi(pageParam)
	pageSize, _ := strconv.Atoi(pageSizeParam)

	tasks, total, err := h.tasksService.GetListPagination(page, pageSize, status, series, sortKey, sortOrder)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp.SuccessResponse(c, http.StatusOK, responses.GetTaskLists{
		Total: total,
		Tasks: tasks,
	})
}

func (h *TasksHandler) DeleteTask(c *gin.Context) {
	var req requests.DeleteTaskRequest
	if err := utils.Bind(c, &req); err != nil {
		return
	}
	var task entities.Task
	_ = copier.Copy(&task, req)

	var resp responses.BaseResponse
	if err := h.tasksService.Delete(&task); err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp.SuccessResponse(c, http.StatusOK, "success")
}

func (h *TasksHandler) GetUploadUrl(c *gin.Context) {
	series := c.Query("series")
	var resp responses.BaseResponse
	if series == "" {
		resp.ErrorResponse(c, http.StatusBadRequest, bizErr.ParseParamsErr)
		return
	}
	url, form, err := h.tasksService.GetUploadUrl(series)
	if err != nil {
		resp.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	resp.SuccessResponse(c, http.StatusOK, responses.GetUploadUrlResponse{
		Url:  url,
		Form: form,
	})
}

func (h *TasksHandler) SetWorkerDevice(c *gin.Context) {
	h.tasksService.SetWorkerDevice(c.Query("device") == "cpu")
	var resp responses.BaseResponse
	resp.SuccessResponse(c, http.StatusOK, "success")
}
