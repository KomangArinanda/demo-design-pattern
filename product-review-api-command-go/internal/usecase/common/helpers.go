package common

import (
	"net/http"

	"example/product-review-api-command-go/internal/shared/appctx"
)

func MustInput[T any](input any) (T, bool) {
	value, ok := input.(T)
	return value, ok
}

func BadRequest(message string) appctx.Response {
	return appctx.Error(http.StatusBadRequest, message)
}
