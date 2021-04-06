package util

import routing "github.com/jackwhelpton/fasthttp-routing"

func DecorateContext(context *routing.Context) *routing.Context {
	context.Response.Header.Set("Access-Control-Allow-Origin", "*")
	return context
}

func ErrorContext(status Status, context *routing.Context, err error) error {
	context.SetStatusCode(status.Code)
	return context.Write(NewMessageObject(err.Error()))
}

func SuccessContext(context *routing.Context, data interface{}) error {
	return context.Write(data)
}

func Throttled(context *routing.Context) error {
	return context.Write(NewMessageObject("You are currently throttled!"))
}

func NewMessageObject(message string) interface{} {
	return struct {
		Message string `json:"message"`
	}{
		Message: message,
	}
}
