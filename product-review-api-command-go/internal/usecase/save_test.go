package usecase_test

import (
	"context"
	"example/product-review-api-command-go/internal/dto/request"
	"example/product-review-api-command-go/internal/dto/response"
	"example/product-review-api-command-go/internal/model"
	reviewuc "example/product-review-api-command-go/internal/usecase/product_review"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSave_Execute(t *testing.T) {
	repository := newFakeRepo(nil)
	usecase := reviewuc.NewSave(repository, fakeOrderClient{
		response: response.OrderValidationResponse{Valid: true, OrderStatus: "COMPLETED"},
	})
	requestModel := reviewuc.SaveRequest{
		ProductID: "PROD-001",
		Request: request.CreateReviewRequest{
			CustomerID: "CUST-001",
			OrderID:    "ORD-001",
			Rating:     5,
			Comment:    "Great",
		},
	}

	result := usecase.Execute(context.Background(), requestModel)

	require.Equal(t, http.StatusCreated, result.Code)
	review, ok := result.Data.(response.CreateReviewResponse)
	require.True(t, ok)
	assert.Equal(t, int64(1), review.ReviewID)
	assert.Equal(t, "PROD-001", review.ProductID)
}

func TestSave_RejectsDuplicate(t *testing.T) {
	repository := newFakeRepo([]model.ProductReview{
		review(1, "PROD-001", "CUST-001", "ORD-001", 5, time.Now()),
	})
	usecase := reviewuc.NewSave(repository, fakeOrderClient{
		response: response.OrderValidationResponse{Valid: true, OrderStatus: "COMPLETED"},
	})
	requestModel := reviewuc.SaveRequest{
		ProductID: "PROD-001",
		Request: request.CreateReviewRequest{
			CustomerID: "CUST-001",
			OrderID:    "ORD-001",
			Rating:     5,
			Comment:    "Great",
		},
	}

	result := usecase.Execute(context.Background(), requestModel)

	require.Equal(t, http.StatusConflict, result.Code)
	assert.Equal(t, reviewuc.ErrDuplicateReview.Error(), result.Message)
}
