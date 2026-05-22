package get_daily_analytics

import (
	"context"
	"example/product-review-api-command-go/internal/dto/response"
	"example/product-review-api-command-go/internal/model"
	"example/product-review-api-command-go/internal/service"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Execute_Success(t *testing.T) {
	repository := newFakeRepo([]model.ProductReview{
		newReview(1, "PROD-001", "CUST-001", "ORD-001", 5, time.Date(2026, 5, 1, 9, 0, 0, 0, time.Local)),
		newReview(2, "PROD-001", "CUST-002", "ORD-002", 3, time.Date(2026, 5, 1, 12, 0, 0, 0, time.Local)),
		newReview(3, "PROD-001", "CUST-003", "ORD-003", 1, time.Date(2026, 5, 2, 8, 0, 0, 0, time.Local)),
	})
	usecase := NewGetDailyAnalytics(repository, service.NewProductReviewService())
	requestModel := GetDailyAnalyticsRequest{ProductID: "PROD-001", Month: 5, Year: 2026}

	result := usecase.Execute(context.Background(), requestModel)

	require.Equal(t, http.StatusOK, result.Code)
	analytics, ok := result.Data.(response.DailyProductReviewAnalyticsResponse)
	require.True(t, ok)
	assert.Equal(t, 3, analytics.TotalReviews)
	assert.Equal(t, 2, len(analytics.DailySummaries))
}

func Test_Execute_RejectsInvalidMonth(t *testing.T) {
	usecase := NewGetDailyAnalytics(newFakeRepo(nil), service.NewProductReviewService())

	result := usecase.Execute(context.Background(), GetDailyAnalyticsRequest{ProductID: "PROD-001", Month: 13, Year: 2026})

	require.Equal(t, http.StatusBadRequest, result.Code)
	assert.Equal(t, ErrInvalidMonth.Error(), result.Message)
}

// -- test doubles --

type fakeRepo struct {
	reviews []model.ProductReview
	nextID  int64
}

func newFakeRepo(reviews []model.ProductReview) *fakeRepo {
	return &fakeRepo{reviews: reviews, nextID: int64(len(reviews) + 1)}
}

func (r *fakeRepo) ExistsByProductCustomerOrder(productID string, customerID string, orderID string) bool {
	return false
}

func (r *fakeRepo) Save(review model.ProductReview) model.ProductReview {
	return model.ProductReview{}
}

func (r *fakeRepo) FindByProductID(productID string) []model.ProductReview {
	return nil
}

func (r *fakeRepo) FindTop5ByProductID(productID string) []model.ProductReview {
	return nil
}

func (r *fakeRepo) FindByProductIDs(productIDs []string) []model.ProductReview {
	return nil
}

func (r *fakeRepo) FindByCustomerID(customerID string) []model.ProductReview {
	return nil
}

func (r *fakeRepo) FindByProductIDOrderByCreatedAtDesc(productID string) []model.ProductReview {
	return nil
}

func (r *fakeRepo) FindRecent(limit int) []model.ProductReview {
	return nil
}

func (r *fakeRepo) FindByProductIDAndCreatedAtBetween(productID string, start time.Time, end time.Time) []model.ProductReview {
	return filterReviews(r.reviews, func(review model.ProductReview) bool {
		return review.ProductID == productID && !review.CreatedAt.Before(start) && review.CreatedAt.Before(end)
	})
}

func (r *fakeRepo) FindByProductIDAndCreatedAtBetweenDesc(productID string, start time.Time, end time.Time) []model.ProductReview {
	return r.FindByProductIDAndCreatedAtBetween(productID, start, end)
}

func newReview(id int64, productID string, customerID string, orderID string, rating int, createdAt time.Time) model.ProductReview {
	return model.ProductReview{
		ID:         id,
		ProductID:  productID,
		CustomerID: customerID,
		OrderID:    orderID,
		Rating:     rating,
		CreatedAt:  createdAt,
		UpdatedAt:  createdAt,
	}
}

func filterReviews(reviews []model.ProductReview, predicate func(model.ProductReview) bool) []model.ProductReview {
	results := make([]model.ProductReview, 0)
	for _, current := range reviews {
		if predicate(current) {
			results = append(results, current)
		}
	}
	return results
}
