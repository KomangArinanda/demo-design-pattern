package save

import (
	"context"
	"example/product-review-api-command-go/internal/client"
	"example/product-review-api-command-go/internal/dto/request"
	"example/product-review-api-command-go/internal/dto/response"
	"example/product-review-api-command-go/internal/model"
	"example/product-review-api-command-go/internal/repo"
	"example/product-review-api-command-go/internal/shared/appctx"
	"example/product-review-api-command-go/internal/usecase/common"
	"net/http"
)

type Request struct {
	ProductID string
	Request   request.CreateReviewRequest
}

type save struct {
	repository  repo.ProductReviewRepo
	orderClient client.OrderClient
}

func NewSave(repository repo.ProductReviewRepo, orderClient client.OrderClient) *save {
	return &save{
		repository:  repository,
		orderClient: orderClient,
	}
}

func (u *save) Execute(_ context.Context, input any) appctx.Response {
	request, ok := common.MustInput[Request](input)
	if !ok {
		return common.BadRequest("Invalid request")
	}
	if request.Request.Rating < 1 || request.Request.Rating > 5 {
		return appctx.Error(http.StatusBadRequest, "rating must be between 1 and 5")
	}
	orderValidation := u.orderClient.ValidateOrder(
		request.Request.CustomerID,
		request.Request.OrderID,
		request.ProductID,
	)
	if !orderValidation.Valid || orderValidation.OrderStatus != "COMPLETED" {
		appctx.Error(http.StatusUnprocessableEntity, "Customer did not purchase this product")
	}
	if u.repository.ExistsByProductCustomerOrder(
		request.ProductID,
		request.Request.CustomerID,
		request.Request.OrderID,
	) {
		return appctx.Error(http.StatusConflict, "duplicate review")
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
