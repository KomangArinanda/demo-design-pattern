package product_review

import "example/product-review-api-command-go/internal/dto/request"

type SaveRequest struct {
	ProductID string
	Request   request.CreateReviewRequest
}

type GetSummaryRequest struct {
	ProductID string
}

type GetSellerReviewAnalyticsRequest struct {
	SellerID string
}

type GetListRequest struct {
	ProductID string
}

type GetListCustomerReviewsRequest struct {
	CustomerID string
}

type GetListRecentReviewsRequest struct {
	Limit int
}

type GetDailyAnalyticsRequest struct {
	ProductID string
	Month     int
	Year      int
}
