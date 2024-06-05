package dto

type ApiError struct {
	Message string `json:"message"`
}

func NewApiError(err error) ApiError {
	return ApiError{
		Message: err.Error(),
	}
}
