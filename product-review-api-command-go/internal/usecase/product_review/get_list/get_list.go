package get_list

import (
	"context"
	"example/product-review-api-command-go/internal/repo"
	"example/product-review-api-command-go/internal/service"
	"example/product-review-api-command-go/internal/shared/appctx"
	"example/product-review-api-command-go/internal/usecase/common"
)

type GetListRequest struct {
	ProductID string
}

type getList struct {
	repository repo.ProductReviewRepo
	service    *service.ProductReviewService
}

func NewGetList(repository repo.ProductReviewRepo, sharedService *service.ProductReviewService) *getList {
	return &getList{repository: repository, service: sharedService}
}

func (u *getList) Execute(_ context.Context, input any) appctx.Response {
	request, ok := common.MustInput[GetListRequest](input)
	if !ok {
		return common.BadRequest("Invalid request")
	}
	reviews := u.repository.FindByProductIDOrderByCreatedAtDesc(request.ProductID)
	return appctx.OK(u.service.MapReviewDetailResponses(reviews))
}
