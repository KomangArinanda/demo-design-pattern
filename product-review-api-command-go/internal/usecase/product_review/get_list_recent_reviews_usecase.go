package product_review

import (
	"context"
	"example/product-review-api-command-go/internal/repo"
	"example/product-review-api-command-go/internal/service"
	"example/product-review-api-command-go/internal/shared/appctx"
	"example/product-review-api-command-go/internal/usecase/common"
)

type getListRecentReviewsUsecase struct {
	repository repo.ProductReviewRepo
	service    *service.ProductReviewService
}

func NewGetListRecentReviewsUsecase(repository repo.ProductReviewRepo, sharedService *service.ProductReviewService) *getListRecentReviewsUsecase {
	return &getListRecentReviewsUsecase{repository: repository, service: sharedService}
}

func (u *getListRecentReviewsUsecase) Execute(_ context.Context, input any) appctx.Response {
	request, ok := common.MustInput[GetListRecentReviewsUsecaseRequest](input)
	if !ok {
		return common.BadRequest("Invalid request")
	}
	limit := u.service.NormalizeRecentReviewLimit(request.Limit)
	reviews := u.repository.FindRecent(limit)
	return appctx.OK(u.service.MapReviewDetailResponses(reviews))
}
