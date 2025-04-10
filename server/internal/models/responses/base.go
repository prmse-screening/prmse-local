package responses

import "github.com/gin-gonic/gin"

type BaseResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   any    `json:"data"`
}

func (response *BaseResponse) ErrorResponse(c *gin.Context, code int, data any) {
	response.Code = code
	response.Status = "error"
	response.Data = data
	c.JSON(code, response)
}

func (response *BaseResponse) SuccessResponse(c *gin.Context, code int, data any) {
	response.Code = code
	response.Status = "success"
	response.Data = data
	c.JSON(code, response)
}
