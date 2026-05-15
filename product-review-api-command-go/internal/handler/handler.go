package handler

import product_review "example/product-review-api-command-go/internal/usecase/product_review"

type Handler struct {
	ProductReviewHandler ProductReviewHandler
}

func NewHandler(usecases *product_review.ProductReviewUsecases) *Handler {
	return &Handler{
		ProductReviewHandler: NewProductReviewHandler(usecases),
	}
}
