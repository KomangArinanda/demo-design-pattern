package product_review

import (
	"context"
	"example/product-review-api-command-go/internal/dto/response"
	"example/product-review-api-command-go/internal/repo"
	"example/product-review-api-command-go/internal/service"
	"example/product-review-api-command-go/internal/shared/appctx"
	"example/product-review-api-command-go/internal/usecase/common"
)

type getSummaryUsecase struct {
	repository repo.ProductReviewRepo
	service    *service.ProductReviewService
}

func NewGetSummaryUsecase(repository repo.ProductReviewRepo, sharedService *service.ProductReviewService) *getSummaryUsecase {
	return &getSummaryUsecase{
		repository: repository,
		service:    sharedService,
	}
}

func (u *getSummaryUsecase) Execute(_ context.Context, input any) appctx.Response {
	request, ok := common.MustInput[GetSummaryUsecaseRequest](input)
	if !ok {
		return common.BadRequest("Invalid request")
	}

	reviews := u.repository.FindByProductID(request.ProductID)
	latest := u.repository.FindTop5ByProductID(request.ProductID)
	latestResponses := u.service.MapLatestReviewResponses(latest)

	return appctx.OK(response.ProductReviewSummaryResponse{
		ProductID:          request.ProductID,
		TotalReviews:       len(reviews),
		AverageRating:      u.service.CalculateAverageRating(reviews),
		RatingDistribution: u.service.CalculateRatingDistribution(reviews),
		LatestReviews:      latestResponses,
	})
}
