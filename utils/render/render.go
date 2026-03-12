package render

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/xuri/excelize/v2"
	"quickBillController/utils/render/errmes"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Code    string `json:"code,omitempty"`
	Count   int64  `json:"count,omitempty"`
}

func WriteExcelFile(c *gin.Context, f *excelize.File, fn string) {
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", url.QueryEscape(fn)))
	c.Header("Content-Transfer-Encoding", "binary")
	_ = f.Write(c.Writer)
}
func QuerySuccess(c *gin.Context, cnt int64, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	c.JSON(http.StatusOK, Response{
		Data:    data,
		Count:   cnt,
		Success: true,
	})
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	c.JSON(http.StatusOK, Response{
		Data:    data,
		Message: "success",
		Success: true,
	})
}

func ResponseError(c *gin.Context, _err errmes.ErrorCode, err error) {
	c.JSON(http.StatusOK, Response{
		Success: false,
		Message: err.Error(),
		Code:    string(_err),
	})
}

func TokenCheck(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, Response{
		Success: false,
		Message: "The login token is invalid. ",
		Code:    string(errmes.ErrTokenInvalid),
	})
}

func BindFailed(c *gin.Context, err error) {
	var message string

	// 处理 EOF 错误（请求体为空）
	if errors.Is(err, io.EOF) {
		message = "Request body is required"
	} else if err.Error() == "EOF" {
		message = "Request body is required"
	} else {
		// 处理验证错误
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			var errMessages []string
			for _, fieldError := range validationErrors {
				errMsg := getValidationErrorMessage(fieldError)
				errMessages = append(errMessages, errMsg)
			}
			message = strings.Join(errMessages, "; ")
		} else {
			// 处理 JSON 解析错误
			if strings.Contains(err.Error(), "json:") {
				message = "Invalid JSON format: " + err.Error()
			} else if strings.Contains(err.Error(), "unexpected") {
				message = "Invalid request format: " + err.Error()
			} else {
				message = err.Error()
			}
		}
	}

	c.JSON(http.StatusOK, Response{
		Success: false,
		Message: message,
		Code:    string(errmes.ErrInvalidRequest),
	})
}

// getValidationErrorMessage 将验证错误转换为友好的错误消息
func getValidationErrorMessage(fieldError validator.FieldError) string {
	field := fieldError.Field()
	tag := fieldError.Tag()

	switch tag {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		return field + " must be at least " + fieldError.Param() + " characters"
	case "max":
		return field + " must be at most " + fieldError.Param() + " characters"
	case "len":
		return field + " must be exactly " + fieldError.Param() + " characters"
	case "gte":
		return field + " must be greater than or equal to " + fieldError.Param()
	case "lte":
		return field + " must be less than or equal to " + fieldError.Param()
	case "gt":
		return field + " must be greater than " + fieldError.Param()
	case "lt":
		return field + " must be less than " + fieldError.Param()
	default:
		return field + " is invalid: " + tag
	}
}
