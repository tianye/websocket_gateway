package http_struct

type Message struct {
	Message      string `json:"message"`
	ErrorMessage string `json:"error_message"`
}

type Response struct {
	Code int     `json:"code"`
	Data Message `json:"data"`
}
