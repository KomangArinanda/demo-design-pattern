package usecase_test

import (
	"example/product-review-api-command-go/internal/dto/response"
	"example/product-review-api-command-go/internal/model"
	"sort"
	"time"
)

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
	lookup := map[string]bool{}
	for _, productID := range productIDs {
		lookup[productID] = true
	}
	return filterReviews(r.reviews, func(review model.ProductReview) bool { return lookup[review.ProductID] })
}

func (r *fakeRepo) FindByCustomerID(customerID string) []model.ProductReview {
	return filterReviews(r.reviews, func(review model.ProductReview) bool { return review.CustomerID == customerID })
}

func (r *fakeRepo) FindByProductIDOrderByCreatedAtDesc(productID string) []model.ProductReview {
	reviews := r.FindByProductID(productID)
	sort.Slice(reviews, func(i, j int) bool {
		return reviews[i].CreatedAt.After(reviews[j].CreatedAt)
	})
	return reviews
}

func (r *fakeRepo) FindRecent(limit int) []model.ProductReview {
	reviews := append([]model.ProductReview(nil), r.reviews...)
	sort.Slice(reviews, func(i, j int) bool {
		return reviews[i].CreatedAt.After(reviews[j].CreatedAt)
	})
	if len(reviews) > limit {
		return reviews[:limit]
	}
	return reviews
}

func (r *fakeRepo) FindByProductIDAndCreatedAtBetween(productID string, start time.Time, end time.Time) []model.ProductReview {
	return filterReviews(r.reviews, func(review model.ProductReview) bool {
		return review.ProductID == productID && !review.CreatedAt.Before(start) && review.CreatedAt.Before(end)
	})
}

func (r *fakeRepo) FindByProductIDAndCreatedAtBetweenDesc(productID string, start time.Time, end time.Time) []model.ProductReview {
	return r.FindByProductIDAndCreatedAtBetween(productID, start, end)
}

type fakeOrderClient struct {
	response response.OrderValidationResponse
}

func (c fakeOrderClient) ValidateOrder(customerID string, orderID string, productID string) response.OrderValidationResponse {
	return c.response
}

type fakeProductClient struct {
	response response.SellerProductsResponse
}

func (c fakeProductClient) GetProductIDsBySellerID(sellerID string) response.SellerProductsResponse {
	return c.response
}

func review(id int64, productID string, customerID string, orderID string, rating int, createdAt time.Time) model.ProductReview {
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
