package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/internal/commons/bizErr"
	"server/internal/models/responses"
)

func Bind[T any](c *gin.Context, request T) error {
	if err := c.ShouldBind(request); err != nil {
		var res responses.BaseResponse
		res.ErrorResponse(c, http.StatusBadRequest, bizErr.ParseParamsErr.Error())
		return err
	}
	return nil
}
