package utils

type GenericError struct {
	Message    string
	StatusCode int
}

func (e GenericError) Error() string {
	return e.Message
}
