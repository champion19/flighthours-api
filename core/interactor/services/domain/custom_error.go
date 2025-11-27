package domain

type AppError struct {
	Code       string
	Message    string
	StatusCode int
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code string, statusCode int) *AppError {
	var message string
	if Messages != nil {
		message = Messages.GetMessage(code)
	} else {
		message = code
	}
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}
