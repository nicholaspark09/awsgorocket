package response

type Response[T any] struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       *T     `json:"data"`
	Error      *error `json:"error"`
}
