package get_summary

import (
	"context"
	"example/product-review-api-command-go/internal/dto/response"
	"example/product-review-api-command-go/internal/model"
	"example/product-review-api-command-go/internal/service"
	"net/http"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Execute_Success(t *testing.T) {
	repository := newFakeRepo([]model.ProductReview{
		newReview(1, "PROD-001", "CUST-001", "ORD-001", 5, time.Now()),
		newReview(2, "PROD-001", "CUST-002", "ORD-002", 3, time.Now()),
	})
	usecase := NewGetSummary(repository, service.NewProductReviewService())

	result := usecase.Execute(context.Background(), GetSummaryRequest{ProductID: "PROD-001"})

	require.Equal(t, http.StatusOK, result.Code)
	summary, ok := result.Data.(response.ProductReviewSummaryResponse)
	require.True(t, ok)
	assert.Equal(t, 2, summary.TotalReviews)
	assert.Equal(t, 4.0, summary.AverageRating)
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
	return filterReviews(r.reviews, func(review model.ProductReview) bool { return review.ProductID == productID })
}

func (r *fakeRepo) FindTop5ByProductID(productID string) []model.ProductReview {
	reviews := r.FindByProductIDOrderByCreatedAtDesc(productID)
	if len(reviews) > 5 {
		return reviews[:5]
	}
	return reviews
}

func (r *fakeRepo) FindByProductIDs(productIDs []string) []model.ProductReview {
	return nil
}

func (r *fakeRepo) FindByCustomerID(customerID string) []model.ProductReview {
	return nil
}

func (r *fakeRepo) FindByProductIDOrderByCreatedAtDesc(productID string) []model.ProductReview {
	reviews := r.FindByProductID(productID)
	sort.Slice(reviews, func(i, j int) bool {
		return reviews[i].CreatedAt.After(reviews[j].CreatedAt)
	})
	return reviews
}

func (r *fakeRepo) FindRecent(limit int) []model.ProductReview {
	return nil
}

func (r *fakeRepo) FindByProductIDAndCreatedAtBetween(productID string, start time.Time, end time.Time) []model.ProductReview {
	return nil
}

func (r *fakeRepo) FindByProductIDAndCreatedAtBetweenDesc(productID string, start time.Time, end time.Time) []model.ProductReview {
	return nil
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
