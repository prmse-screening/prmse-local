package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/internal/services"
	"strconv"
)

type DicomHandler struct {
	dicomService *services.DicomService
}

func NewDicomHandler(dicomService *services.DicomService) *DicomHandler {
	return &DicomHandler{
		dicomService: dicomService,
	}
}

func (h *DicomHandler) GetUrl(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return
	}
	url, err := h.dicomService.GetUrl(id)
	if err != nil {
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, url)
}
