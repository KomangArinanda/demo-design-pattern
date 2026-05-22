package get_list_customer_reviews

import (
	"context"
	"example/product-review-api-command-go/internal/repo"
	"example/product-review-api-command-go/internal/service"
	"example/product-review-api-command-go/internal/shared/appctx"
	"example/product-review-api-command-go/internal/usecase/common"
)

type GetListCustomerReviewsRequest struct {
	CustomerID string
}

type getListCustomerReviews struct {
	repository repo.ProductReviewRepo
	service    *service.ProductReviewService
}

func NewGetListCustomerReviews(repository repo.ProductReviewRepo, sharedService *service.ProductReviewService) *getListCustomerReviews {
	return &getListCustomerReviews{repository: repository, service: sharedService}
}

func (u *getListCustomerReviews) Execute(_ context.Context, input any) appctx.Response {
	request, ok := common.MustInput[GetListCustomerReviewsRequest](input)
	if !ok {
		return common.BadRequest("Invalid request")
	}
	reviews := u.repository.FindByCustomerID(request.CustomerID)
	return appctx.OK(u.service.MapReviewDetailResponses(reviews))
}
