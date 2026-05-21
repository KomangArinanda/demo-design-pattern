package handler

import "example/product-review-api-command-go/internal/usecase/product_review"

type Handler struct {
	ProductReviewHandler ProductReviewHandler
}

func NewHandler(usecases *product_review.ProductReview) *Handler {
	return &Handler{
		ProductReviewHandler: NewProductReviewHandler(usecases),
	}
}
