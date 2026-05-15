package product_review

import (
	"context"
	"example/product-review-api-command-go/internal/repo"
	"example/product-review-api-command-go/internal/service"
	"example/product-review-api-command-go/internal/shared/appctx"
	"example/product-review-api-command-go/internal/usecase/common"
)

type getListUsecase struct {
	repository repo.ProductReviewRepo
	service    *service.ProductReviewService
}

func NewGetListUsecase(repository repo.ProductReviewRepo, sharedService *service.ProductReviewService) *getListUsecase {
	return &getListUsecase{repository: repository, service: sharedService}
}

func (u *getListUsecase) Execute(_ context.Context, input any) appctx.Response {
	request, ok := common.MustInput[GetListUsecaseRequest](input)
	if !ok {
		return common.BadRequest("Invalid request")
	}
	reviews := u.repository.FindByProductIDOrderByCreatedAtDesc(request.ProductID)
	return appctx.OK(u.service.MapReviewDetailResponses(reviews))
}
