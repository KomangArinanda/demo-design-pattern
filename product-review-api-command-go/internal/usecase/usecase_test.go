package usecase_test

import (
	"context"
	"example/product-review-api-command-go/internal/dto/request"
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

func TestSaveUsecase_Execute(t *testing.T) {
	repository := newFakeRepo(nil)
	usecase := reviewuc.NewSaveUsecase(repository, fakeOrderClient{
		response: response.OrderValidationResponse{Valid: true, OrderStatus: "COMPLETED"},
	})
	requestModel := reviewuc.SaveUsecaseRequest{
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

func TestSaveUsecase_RejectsDuplicate(t *testing.T) {
	repository := newFakeRepo([]model.ProductReview{
		review(1, "PROD-001", "CUST-001", "ORD-001", 5, time.Now()),
	})
	usecase := reviewuc.NewSaveUsecase(repository, fakeOrderClient{
		response: response.OrderValidationResponse{Valid: true, OrderStatus: "COMPLETED"},
	})
	requestModel := reviewuc.SaveUsecaseRequest{
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

func TestGetSummaryUsecase_Execute(t *testing.T) {
	repository := newFakeRepo([]model.ProductReview{
		review(1, "PROD-001", "CUST-001", "ORD-001", 5, time.Now()),
		review(2, "PROD-001", "CUST-002", "ORD-002", 3, time.Now()),
	})
	usecase := reviewuc.NewGetSummaryUsecase(repository, service.NewProductReviewService())

	result := usecase.Execute(context.Background(), reviewuc.GetSummaryUsecaseRequest{ProductID: "PROD-001"})

	require.Equal(t, http.StatusOK, result.Code)
	summary, ok := result.Data.(response.ProductReviewSummaryResponse)
	require.True(t, ok)
	assert.Equal(t, 2, summary.TotalReviews)
	assert.Equal(t, 4.0, summary.AverageRating)
}

func TestGetSellerReviewAnalyticsUsecase_Execute(t *testing.T) {
	repository := newFakeRepo([]model.ProductReview{
		review(1, "PROD-001", "CUST-001", "ORD-001", 5, time.Now()),
		review(2, "PROD-002", "CUST-002", "ORD-002", 2, time.Now()),
		review(3, "PROD-003", "CUST-003", "ORD-003", 1, time.Now()),
	})
	usecase := reviewuc.NewGetSellerReviewAnalyticsUsecase(repository, fakeProductClient{
		response: response.SellerProductsResponse{SellerID: "SELLER-001", ProductIDs: []string{"PROD-001", "PROD-002", "PROD-003"}},
	}, service.NewProductReviewService())

	result := usecase.Execute(context.Background(), reviewuc.GetSellerReviewAnalyticsUsecaseRequest{SellerID: "SELLER-001"})

	require.Equal(t, http.StatusOK, result.Code)
	analytics, ok := result.Data.(response.SellerReviewAnalyticsResponse)
	require.True(t, ok)
	assert.Equal(t, 3, analytics.TotalReviews)
	assert.Equal(t, 2, analytics.NegativeReviewCount)
}

func TestGetListUsecase_Execute(t *testing.T) {
	repository := newFakeRepo([]model.ProductReview{
		review(1, "PROD-001", "CUST-001", "ORD-001", 5, time.Date(2026, 5, 2, 9, 0, 0, 0, time.Local)),
		review(2, "PROD-001", "CUST-002", "ORD-002", 4, time.Date(2026, 5, 3, 9, 0, 0, 0, time.Local)),
	})
	usecase := reviewuc.NewGetListUsecase(repository, service.NewProductReviewService())

	result := usecase.Execute(context.Background(), reviewuc.GetListUsecaseRequest{ProductID: "PROD-001"})

	require.Equal(t, http.StatusOK, result.Code)
	reviews, ok := result.Data.([]response.ReviewDetailResponse)
	require.True(t, ok)
	require.Len(t, reviews, 2)
	assert.Equal(t, int64(2), reviews[0].ReviewID)
}

func TestGetListCustomerReviewsUsecase_Execute(t *testing.T) {
	repository := newFakeRepo([]model.ProductReview{
		review(1, "PROD-001", "CUST-001", "ORD-001", 5, time.Now()),
		review(2, "PROD-002", "CUST-001", "ORD-002", 4, time.Now()),
	})
	usecase := reviewuc.NewGetListCustomerReviewsUsecase(repository, service.NewProductReviewService())

	result := usecase.Execute(context.Background(), reviewuc.GetListCustomerReviewsUsecaseRequest{CustomerID: "CUST-001"})

	require.Equal(t, http.StatusOK, result.Code)
	reviews, ok := result.Data.([]response.ReviewDetailResponse)
	require.True(t, ok)
	require.Len(t, reviews, 2)
}

func TestGetListRecentReviewsUsecase_Execute(t *testing.T) {
	repository := newFakeRepo([]model.ProductReview{
		review(1, "PROD-001", "CUST-001", "ORD-001", 5, time.Date(2026, 5, 1, 9, 0, 0, 0, time.Local)),
		review(2, "PROD-002", "CUST-002", "ORD-002", 4, time.Date(2026, 5, 3, 9, 0, 0, 0, time.Local)),
		review(3, "PROD-003", "CUST-003", "ORD-003", 3, time.Date(2026, 5, 2, 9, 0, 0, 0, time.Local)),
	})
	usecase := reviewuc.NewGetListRecentReviewsUsecase(repository, service.NewProductReviewService())

	result := usecase.Execute(context.Background(), reviewuc.GetListRecentReviewsUsecaseRequest{Limit: 2})

	require.Equal(t, http.StatusOK, result.Code)
	reviews, ok := result.Data.([]response.ReviewDetailResponse)
	require.True(t, ok)
	require.Len(t, reviews, 2)
	assert.Equal(t, int64(2), reviews[0].ReviewID)
}

func TestGetDailyProductReviewAnalyticsUsecase_Execute(t *testing.T) {
	repository := newFakeRepo([]model.ProductReview{
		review(1, "PROD-001", "CUST-001", "ORD-001", 5, time.Date(2026, 5, 1, 9, 0, 0, 0, time.Local)),
		review(2, "PROD-001", "CUST-002", "ORD-002", 3, time.Date(2026, 5, 1, 12, 0, 0, 0, time.Local)),
		review(3, "PROD-001", "CUST-003", "ORD-003", 1, time.Date(2026, 5, 2, 8, 0, 0, 0, time.Local)),
	})
	usecase := reviewuc.NewGetDailyAnalyticsUsecase(repository, service.NewProductReviewService())
	requestModel := reviewuc.GetDailyAnalyticsUsecaseRequest{ProductID: "PROD-001", Month: 5, Year: 2026}

	result := usecase.Execute(context.Background(), requestModel)

	require.Equal(t, http.StatusOK, result.Code)
	analytics, ok := result.Data.(response.DailyProductReviewAnalyticsResponse)
	require.True(t, ok)
	assert.Equal(t, 3, analytics.TotalReviews)
	assert.Equal(t, 2, len(analytics.DailySummaries))
}

func TestGetDailyProductReviewAnalyticsUsecase_RejectsInvalidMonth(t *testing.T) {
	usecase := reviewuc.NewGetDailyAnalyticsUsecase(newFakeRepo(nil), service.NewProductReviewService())

	result := usecase.Execute(context.Background(), reviewuc.GetDailyAnalyticsUsecaseRequest{ProductID: "PROD-001", Month: 13, Year: 2026})

	require.Equal(t, http.StatusBadRequest, result.Code)
	assert.Equal(t, reviewuc.ErrInvalidMonth.Error(), result.Message)
}
