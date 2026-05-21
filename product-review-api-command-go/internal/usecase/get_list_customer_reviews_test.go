package usecase_test

import (
	"context"
	"example/product-review-api-command-go/internal/dto/response"
	"example/product-review-api-command-go/internal/model"
	"example/product-review-api-command-go/internal/service"
	reviewuc "example/product-review-api-command-go/internal/usecase/product_review"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetListCustomerReviews_Execute(t *testing.T) {
	repository := newFakeRepo([]model.ProductReview{
		review(1, "PROD-001", "CUST-001", "ORD-001", 5, time.Now()),
		review(2, "PROD-002", "CUST-001", "ORD-002", 4, time.Now()),
	})
	usecase := reviewuc.NewGetListCustomerReviews(repository, service.NewProductReviewService())

	result := usecase.Execute(context.Background(), reviewuc.GetListCustomerReviewsRequest{CustomerID: "CUST-001"})

	require.Equal(t, http.StatusOK, result.Code)
	reviews, ok := result.Data.([]response.ReviewDetailResponse)
	require.True(t, ok)
	require.Len(t, reviews, 2)
}
