package response

type NetworkResponse[T any] struct {
	Data       *T
	StatusCode int
	Message    *string
}

func NewSuccessResponse[T any](data *T) *NetworkResponse[T] {
	return &NetworkResponse[T]{Data: data, StatusCode: 200}
}

func NewErrorResponse[T any](data *T, code *int, reason *string) *NetworkResponse[T] {
	var statusCode int = 500
	if code != nil {
		statusCode = *code
	}
	return &NetworkResponse[T]{Data: data, StatusCode: statusCode, Message: reason}
}
