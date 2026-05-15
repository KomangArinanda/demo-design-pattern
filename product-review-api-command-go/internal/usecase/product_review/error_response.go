package product_review

import (
	"net/http"

	"example/product-review-api-command-go/internal/shared/appctx"
)

func errorResponse(err error) appctx.Response {
	switch err {
	case ErrDuplicateReview:
		return appctx.Error(http.StatusConflict, err.Error())
	case ErrInvalidOrder:
		return appctx.Error(http.StatusUnprocessableEntity, err.Error())
	case ErrInvalidRating, ErrInvalidMonth:
		return appctx.Error(http.StatusBadRequest, err.Error())
	default:
		return appctx.Error(http.StatusBadRequest, err.Error())
	}
}
