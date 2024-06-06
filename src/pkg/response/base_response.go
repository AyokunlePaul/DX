package response

import "net/http"

type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Code    int         `json:"-"`
}

func NewCreatedResponse(message string, data interface{}) *BaseResponse {
	return &BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
		Code:    http.StatusCreated,
	}
}

func NewOkResponse(message string, data interface{}) *BaseResponse {
	return &BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
		Code:    http.StatusOK,
	}
}

func NewNotFoundError(message string) *BaseResponse {
	return &BaseResponse{
		Success: false,
		Message: message,
		Code:    http.StatusNotFound,
		Data:    nil,
	}
}

func NewBadRequestError(message string) *BaseResponse {
	return &BaseResponse{
		Success: false,
		Code:    http.StatusBadRequest,
		Message: message,
		Data:    nil,
	}
}

func NewInternalServerError(message string) *BaseResponse {
	return &BaseResponse{
		Success: false,
		Code:    http.StatusInternalServerError,
		Message: message,
		Data:    nil,
	}
}

func NewUnAuthorizedError() *BaseResponse {
	return &BaseResponse{
		Success: false,
		Code:    http.StatusUnauthorized,
		Message: Unauthorized,
		Data:    nil,
	}
}
