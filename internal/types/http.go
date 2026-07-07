package types

type ResponseHTTP struct {
	Data any `json:"data,omitempty"`
	Code int `json:"code"`
}

type ResponseError struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Code    int    `json:"code"`
}
