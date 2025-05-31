package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponstData struct {
	Code    int         `json:"code"`           //status code
	Message string      `json:"message"`        // thong bao loi
	Data    interface{} `json:"data,omitempty"` // du lieu tra ve
}

// success response

func SuccessResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, ResponstData{
		Code:    code,
		Message: msg[code],
		Data:    data,
	})
}


// error response
func ErrorResponse(c *gin.Context, code int, message string) {
	if message == "" {
		message = msg[code]
	}
	c.JSON(http.StatusOK, ResponstData{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}