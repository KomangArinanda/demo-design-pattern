package product_review

import (
	"example/product-review-api-command-go/internal/client"
	"example/product-review-api-command-go/internal/repo"
	"example/product-review-api-command-go/internal/service"
	"example/product-review-api-command-go/internal/usecase"
)

type ProductReviewUsecases struct {
	Save                           usecase.Usecase
	GetSummary                     usecase.Usecase
	GetSellerReviewAnalytics       usecase.Usecase
	GetList                        usecase.Usecase
	GetListCustomerReviews         usecase.Usecase
	GetListRecentReviews           usecase.Usecase
	GetDailyProductReviewAnalytics usecase.Usecase
}

func NewProductReviewUsecases(
	repository repo.ProductReviewRepo,
	orderClient client.OrderClient,
	productClient client.ProductClient,
	sharedService *service.ProductReviewService,
) *ProductReviewUsecases {
	return &ProductReviewUsecases{
		Save:                           NewSaveUsecase(repository, orderClient),
		GetSummary:                     NewGetSummaryUsecase(repository, sharedService),
		GetSellerReviewAnalytics:       NewGetSellerReviewAnalyticsUsecase(repository, productClient, sharedService),
		GetList:                        NewGetListUsecase(repository, sharedService),
		GetListCustomerReviews:         NewGetListCustomerReviewsUsecase(repository, sharedService),
		GetListRecentReviews:           NewGetListRecentReviewsUsecase(repository, sharedService),
		GetDailyProductReviewAnalytics: NewGetDailyAnalyticsUsecase(repository, sharedService),
	}
}
