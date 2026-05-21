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

func TestGetSummary_Execute(t *testing.T) {
	repository := newFakeRepo([]model.ProductReview{
		review(1, "PROD-001", "CUST-001", "ORD-001", 5, time.Now()),
		review(2, "PROD-001", "CUST-002", "ORD-002", 3, time.Now()),
	})
	usecase := reviewuc.NewGetSummary(repository, service.NewProductReviewService())

	result := usecase.Execute(context.Background(), reviewuc.GetSummaryRequest{ProductID: "PROD-001"})

	require.Equal(t, http.StatusOK, result.Code)
	summary, ok := result.Data.(response.ProductReviewSummaryResponse)
	require.True(t, ok)
	assert.Equal(t, 2, summary.TotalReviews)
	assert.Equal(t, 4.0, summary.AverageRating)
}
