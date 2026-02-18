package util

type AppError struct {
	Error   error
	Code    int
	Message string
}

func Error(err error, code int, message string) *AppError {
	return &AppError{Error: err, Code: code, Message: message}
}
