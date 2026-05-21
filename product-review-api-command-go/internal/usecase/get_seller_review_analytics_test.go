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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSellerReviewAnalytics_Execute(t *testing.T) {
	repository := newFakeRepo([]model.ProductReview{
		review(1, "PROD-001", "CUST-001", "ORD-001", 5, time.Now()),
		review(2, "PROD-002", "CUST-002", "ORD-002", 2, time.Now()),
		review(3, "PROD-003", "CUST-003", "ORD-003", 1, time.Now()),
	})
	usecase := reviewuc.NewGetSellerReviewAnalytics(repository, fakeProductClient{
		response: response.SellerProductsResponse{SellerID: "SELLER-001", ProductIDs: []string{"PROD-001", "PROD-002", "PROD-003"}},
	}, service.NewProductReviewService())

	result := usecase.Execute(context.Background(), reviewuc.GetSellerReviewAnalyticsRequest{SellerID: "SELLER-001"})

	require.Equal(t, http.StatusOK, result.Code)
	analytics, ok := result.Data.(response.SellerReviewAnalyticsResponse)
	require.True(t, ok)
	assert.Equal(t, 3, analytics.TotalReviews)
	assert.Equal(t, 2, analytics.NegativeReviewCount)
}
