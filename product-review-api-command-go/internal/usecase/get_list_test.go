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

func TestGetList_Execute(t *testing.T) {
	repository := newFakeRepo([]model.ProductReview{
		review(1, "PROD-001", "CUST-001", "ORD-001", 5, time.Date(2026, 5, 2, 9, 0, 0, 0, time.Local)),
		review(2, "PROD-001", "CUST-002", "ORD-002", 4, time.Date(2026, 5, 3, 9, 0, 0, 0, time.Local)),
	})
	usecase := reviewuc.NewGetList(repository, service.NewProductReviewService())

	result := usecase.Execute(context.Background(), reviewuc.GetListRequest{ProductID: "PROD-001"})

	require.Equal(t, http.StatusOK, result.Code)
	reviews, ok := result.Data.([]response.ReviewDetailResponse)
	require.True(t, ok)
	require.Len(t, reviews, 2)
	assert.Equal(t, int64(2), reviews[0].ReviewID)
}
