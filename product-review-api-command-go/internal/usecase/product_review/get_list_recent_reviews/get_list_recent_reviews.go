package get_list_recent_reviews

import (
	"context"
	"example/product-review-api-command-go/internal/repo"
	"example/product-review-api-command-go/internal/service"
	"example/product-review-api-command-go/internal/shared/appctx"
	"example/product-review-api-command-go/internal/usecase/common"
)

type GetListRecentReviewsRequest struct {
	Limit int
}

type getListRecentReviews struct {
	repository repo.ProductReviewRepo
	service    *service.ProductReviewService
}

func NewGetListRecentReviews(repository repo.ProductReviewRepo, sharedService *service.ProductReviewService) *getListRecentReviews {
	return &getListRecentReviews{repository: repository, service: sharedService}
}

func (u *getListRecentReviews) Execute(_ context.Context, input any) appctx.Response {
	request, ok := common.MustInput[GetListRecentReviewsRequest](input)
	if !ok {
		return common.BadRequest("Invalid request")
	}
	limit := u.service.NormalizeRecentReviewLimit(request.Limit)
	reviews := u.repository.FindRecent(limit)
	return appctx.OK(u.service.MapReviewDetailResponses(reviews))
}
