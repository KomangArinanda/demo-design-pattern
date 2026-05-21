package product_review

import (
	"example/product-review-api-command-go/internal/client"
	"example/product-review-api-command-go/internal/repo"
	"example/product-review-api-command-go/internal/service"
	"example/product-review-api-command-go/internal/usecase"
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
		Save:                           NewSave(repository, orderClient),
		GetSummary:                     NewGetSummary(repository, sharedService),
		GetSellerReviewAnalytics:       NewGetSellerReviewAnalytics(repository, productClient, sharedService),
		GetList:                        NewGetList(repository, sharedService),
		GetListCustomerReviews:         NewGetListCustomerReviews(repository, sharedService),
		GetListRecentReviews:           NewGetListRecentReviews(repository, sharedService),
		GetDailyProductReviewAnalytics: NewGetDailyAnalytics(repository, sharedService),
	}
}
