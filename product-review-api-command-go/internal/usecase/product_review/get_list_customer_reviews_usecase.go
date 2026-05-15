package product_review

import (
	"context"
	"example/product-review-api-command-go/internal/repo"
	"example/product-review-api-command-go/internal/service"
	"example/product-review-api-command-go/internal/shared/appctx"
	"example/product-review-api-command-go/internal/usecase/common"
)

type getListCustomerReviewsUsecase struct {
	repository repo.ProductReviewRepo
	service    *service.ProductReviewService
}

func NewGetListCustomerReviewsUsecase(repository repo.ProductReviewRepo, sharedService *service.ProductReviewService) *getListCustomerReviewsUsecase {
	return &getListCustomerReviewsUsecase{repository: repository, service: sharedService}
}

func (u *getListCustomerReviewsUsecase) Execute(_ context.Context, input any) appctx.Response {
	request, ok := common.MustInput[GetListCustomerReviewsUsecaseRequest](input)
	if !ok {
		return common.BadRequest("Invalid request")
	}
	reviews := u.repository.FindByCustomerID(request.CustomerID)
	return appctx.OK(u.service.MapReviewDetailResponses(reviews))
}
