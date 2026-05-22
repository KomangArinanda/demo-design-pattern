package product_review

import (
	"example/product-review-api-command-go/internal/client"
	"example/product-review-api-command-go/internal/repo"
	"example/product-review-api-command-go/internal/service"
	"example/product-review-api-command-go/internal/usecase"
	"example/product-review-api-command-go/internal/usecase/product_review/get_daily_analytics"
	"example/product-review-api-command-go/internal/usecase/product_review/get_list"
	"example/product-review-api-command-go/internal/usecase/product_review/get_list_customer_reviews"
	"example/product-review-api-command-go/internal/usecase/product_review/get_list_recent_reviews"
	"example/product-review-api-command-go/internal/usecase/product_review/get_seller_review_analytics"
	"example/product-review-api-command-go/internal/usecase/product_review/get_summary"
	"example/product-review-api-command-go/internal/usecase/product_review/save"
)

type ProductReview struct {
	Save                           usecase.Executor
	GetSummary                     usecase.Executor
	GetSellerReviewAnalytics       usecase.Executor
	GetList                        usecase.Executor
	GetListCustomerReviews         usecase.Executor
	GetListRecentReviews           usecase.Executor
	GetDailyProductReviewAnalytics usecase.Executor
}

func NewProductReview(
	repository repo.ProductReviewRepo,
	orderClient client.OrderClient,
	productClient client.ProductClient,
	sharedService *service.ProductReviewService,
) *ProductReview {
	return &ProductReview{
		Save:                           save.NewSave(repository, orderClient),
		GetSummary:                     get_summary.NewGetSummary(repository, sharedService),
		GetSellerReviewAnalytics:       get_seller_review_analytics.NewGetSellerReviewAnalytics(repository, productClient, sharedService),
		GetList:                        get_list.NewGetList(repository, sharedService),
		GetListCustomerReviews:         get_list_customer_reviews.NewGetListCustomerReviews(repository, sharedService),
		GetListRecentReviews:           get_list_recent_reviews.NewGetListRecentReviews(repository, sharedService),
		GetDailyProductReviewAnalytics: get_daily_analytics.NewGetDailyAnalytics(repository, sharedService),
	}
}
