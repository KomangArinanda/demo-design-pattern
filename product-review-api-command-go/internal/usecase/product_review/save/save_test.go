package save

import (
	"context"
	"example/product-review-api-command-go/internal/dto/request"
	"example/product-review-api-command-go/internal/dto/response"
	"example/product-review-api-command-go/internal/model"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Execute_Success(t *testing.T) {
	repository := newFakeRepo(nil)
	usecase := NewSave(repository, fakeOrderClient{
		response: response.OrderValidationResponse{Valid: true, OrderStatus: "COMPLETED"},
	})
	requestModel := Request{
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

func Test_Execute_RejectsDuplicate(t *testing.T) {
	repository := newFakeRepo([]model.ProductReview{
		newReview(1, "PROD-001", "CUST-001", "ORD-001", 5, time.Now()),
	})
	usecase := NewSave(repository, fakeOrderClient{
		response: response.OrderValidationResponse{Valid: true, OrderStatus: "COMPLETED"},
	})
	requestModel := Request{
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
	assert.Equal(t, "duplicate review", result.Message)
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
	for _, review := range r.reviews {
		if review.ProductID == productID && review.CustomerID == customerID && review.OrderID == orderID {
			return true
		}
	}
	return false
}

func (r *fakeRepo) Save(review model.ProductReview) model.ProductReview {
	review.ID = r.nextID
	review.CreatedAt = time.Date(2026, 5, 15, 10, 0, 0, 0, time.Local)
	review.UpdatedAt = review.CreatedAt
	r.nextID++
	r.reviews = append(r.reviews, review)
	return review
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
	return nil
}

func (r *fakeRepo) FindByProductIDAndCreatedAtBetweenDesc(productID string, start time.Time, end time.Time) []model.ProductReview {
	return nil
}

type fakeOrderClient struct {
	response response.OrderValidationResponse
}

func (c fakeOrderClient) ValidateOrder(customerID string, orderID string, productID string) response.OrderValidationResponse {
	return c.response
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
