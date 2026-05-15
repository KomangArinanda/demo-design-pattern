package appctx

import "net/http"

type Response struct {
	Code    int
	Message string
	Data    any
}

func OK(data any) Response {
	return Response{Code: http.StatusOK, Data: data}
}

func Created(data any) Response {
	return Response{Code: http.StatusCreated, Data: data}
}

func Error(code int, message string) Response {
	return Response{Code: code, Message: message}
}
