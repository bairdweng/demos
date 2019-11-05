package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 响应对象
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// JSON 统一返回格式
func JSON(c *gin.Context, code int, message string, data ...interface{}) {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}

	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    responseData,
	})

}
