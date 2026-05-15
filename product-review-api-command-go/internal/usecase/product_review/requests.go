package product_review

import "example/product-review-api-command-go/internal/dto/request"

type SaveUsecaseRequest struct {
	ProductID string
	Request   request.CreateReviewRequest
}

type GetSummaryUsecaseRequest struct {
	ProductID string
}

type GetSellerReviewAnalyticsUsecaseRequest struct {
	SellerID string
}

type GetListUsecaseRequest struct {
	ProductID string
}

type GetListCustomerReviewsUsecaseRequest struct {
	CustomerID string
}

type GetListRecentReviewsUsecaseRequest struct {
	Limit int
}

type GetDailyAnalyticsUsecaseRequest struct {
	ProductID string
	Month     int
	Year      int
}
