package dto

type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func SuccessResponse(message string, data interface{}) BaseResponse {
	return BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(code, message string, details map[string]interface{}) *BaseResponse {
	return &BaseResponse{
		Success: false,
		Message: message,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

func ValidationErrorResponse(validationErrors map[string]string) *BaseResponse {
	details := make(map[string]interface{})
	details["validation_errors"] = validationErrors

	return &BaseResponse{
		Success: false,
		Message: "Validation failed",
		Error: &ErrorInfo{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid input data",
			Details: details,
		},
	}
}
