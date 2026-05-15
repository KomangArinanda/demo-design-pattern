package product_review

import (
	"context"
	"example/product-review-api-command-go/internal/client"
	"example/product-review-api-command-go/internal/dto/response"
	"example/product-review-api-command-go/internal/model"
	"example/product-review-api-command-go/internal/repo"
	"example/product-review-api-command-go/internal/shared/appctx"
	"example/product-review-api-command-go/internal/usecase/common"
)

type saveUsecase struct {
	repository  repo.ProductReviewRepo
	orderClient client.OrderClient
}

func NewSaveUsecase(repository repo.ProductReviewRepo, orderClient client.OrderClient) *saveUsecase {
	return &saveUsecase{
		repository:  repository,
		orderClient: orderClient,
	}
}

func (u *saveUsecase) Execute(_ context.Context, input any) appctx.Response {
	request, ok := common.MustInput[SaveUsecaseRequest](input)
	if !ok {
		return common.BadRequest("Invalid request")
	}
	if request.Request.Rating < 1 || request.Request.Rating > 5 {
		return errorResponse(ErrInvalidRating)
	}
	orderValidation := u.orderClient.ValidateOrder(
		request.Request.CustomerID,
		request.Request.OrderID,
		request.ProductID,
	)
	if !orderValidation.Valid || orderValidation.OrderStatus != "COMPLETED" {
		return errorResponse(ErrInvalidOrder)
	}
	if u.repository.ExistsByProductCustomerOrder(
		request.ProductID,
		request.Request.CustomerID,
		request.Request.OrderID,
	) {
		return errorResponse(ErrDuplicateReview)
	}

	saved := u.repository.Save(model.ProductReview{
		ProductID:  request.ProductID,
		CustomerID: request.Request.CustomerID,
		OrderID:    request.Request.OrderID,
		Rating:     request.Request.Rating,
		Comment:    request.Request.Comment,
	})

	return appctx.Created(response.CreateReviewResponse{
		ReviewID:   saved.ID,
		ProductID:  saved.ProductID,
		CustomerID: saved.CustomerID,
		OrderID:    saved.OrderID,
		Rating:     saved.Rating,
		Comment:    saved.Comment,
		CreatedAt:  saved.CreatedAt,
	})
}
