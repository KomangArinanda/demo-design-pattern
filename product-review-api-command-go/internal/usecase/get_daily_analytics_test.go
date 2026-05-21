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

func TestGetDailyProductReviewAnalytics_Execute(t *testing.T) {
	repository := newFakeRepo([]model.ProductReview{
		review(1, "PROD-001", "CUST-001", "ORD-001", 5, time.Date(2026, 5, 1, 9, 0, 0, 0, time.Local)),
		review(2, "PROD-001", "CUST-002", "ORD-002", 3, time.Date(2026, 5, 1, 12, 0, 0, 0, time.Local)),
		review(3, "PROD-001", "CUST-003", "ORD-003", 1, time.Date(2026, 5, 2, 8, 0, 0, 0, time.Local)),
	})
	usecase := reviewuc.NewGetDailyAnalytics(repository, service.NewProductReviewService())
	requestModel := reviewuc.GetDailyAnalyticsRequest{ProductID: "PROD-001", Month: 5, Year: 2026}

	result := usecase.Execute(context.Background(), requestModel)

	require.Equal(t, http.StatusOK, result.Code)
	analytics, ok := result.Data.(response.DailyProductReviewAnalyticsResponse)
	require.True(t, ok)
	assert.Equal(t, 3, analytics.TotalReviews)
	assert.Equal(t, 2, len(analytics.DailySummaries))
}

func TestGetDailyProductReviewAnalytics_RejectsInvalidMonth(t *testing.T) {
	usecase := reviewuc.NewGetDailyAnalytics(newFakeRepo(nil), service.NewProductReviewService())

	result := usecase.Execute(context.Background(), reviewuc.GetDailyAnalyticsRequest{ProductID: "PROD-001", Month: 13, Year: 2026})

	require.Equal(t, http.StatusBadRequest, result.Code)
	assert.Equal(t, reviewuc.ErrInvalidMonth.Error(), result.Message)
}
