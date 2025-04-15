package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/internal/services"
)

type DicomHandler struct {
	dicomService *services.DicomService
}

func NewDicomHandler(dicomService *services.DicomService) *DicomHandler {
	return &DicomHandler{
		dicomService: dicomService,
	}
}

func (h *DicomHandler) Redirect(c *gin.Context) {
	series := c.Param("series")
	file := c.Param("file")
	url, err := h.dicomService.Redirect(series, file)
	if err != nil {
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, url)
}
