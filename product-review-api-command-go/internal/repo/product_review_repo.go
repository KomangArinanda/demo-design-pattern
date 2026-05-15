package repo

import (
	"example/product-review-api-command-go/internal/model"
	"sort"
	"strconv"
	"sync"
	"time"
)

type ProductReviewRepo interface {
	ExistsByProductCustomerOrder(productID string, customerID string, orderID string) bool
	Save(review model.ProductReview) model.ProductReview
	FindByProductID(productID string) []model.ProductReview
	FindTop5ByProductID(productID string) []model.ProductReview
	FindByProductIDs(productIDs []string) []model.ProductReview
	FindByCustomerID(customerID string) []model.ProductReview
	FindByProductIDOrderByCreatedAtDesc(productID string) []model.ProductReview
	FindRecent(limit int) []model.ProductReview
	FindByProductIDAndCreatedAtBetween(productID string, start time.Time, end time.Time) []model.ProductReview
	FindByProductIDAndCreatedAtBetweenDesc(productID string, start time.Time, end time.Time) []model.ProductReview
}

type productReviewRepo struct {
	mu      sync.RWMutex
	nextID  int64
	reviews []model.ProductReview
	latency time.Duration
}

func NewProductReviewRepo(latency time.Duration) *productReviewRepo {
	return &productReviewRepo{
		nextID:  1,
		reviews: make([]model.ProductReview, 0),
		latency: latency,
	}
}

func (r *productReviewRepo) Seed() {
	r.seedProduct("PROD-001", []int{5, 5, 4, 4, 5, 3, 4, 5, 2, 5})
	r.seedProduct("PROD-002", []int{5, 4, 4, 3, 2, 2, 1, 4})
	r.seedProduct("PROD-003", []int{5, 5, 4, 2, 1})
	r.seedProduct("PROD-004", []int{5, 4, 4, 4, 3, 2, 5})
	r.seedProduct("PROD-005", []int{5, 4, 3, 1})
}

func (r *productReviewRepo) seedProduct(productID string, ratings []int) {
	for index, rating := range ratings {
		suffix := productID[len(productID)-3:]
		r.Save(model.ProductReview{
			ProductID:  productID,
			CustomerID: "CUST-" + suffix + "-" + itoa(index+1),
			OrderID:    "ORD-" + suffix + "-" + itoa(index+1),
			Rating:     rating,
			Comment:    "Sample review " + itoa(index+1) + " for " + productID,
		})
	}
}

func (r *productReviewRepo) delay() {
	if r.latency > 0 {
		time.Sleep(r.latency)
	}
}

func (r *productReviewRepo) ExistsByProductCustomerOrder(productID string, customerID string, orderID string) bool {
	r.delay()
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, review := range r.reviews {
		if review.ProductID == productID && review.CustomerID == customerID && review.OrderID == orderID {
			return true
		}
	}
	return false
}

func (r *productReviewRepo) Save(review model.ProductReview) model.ProductReview {
	r.delay()
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	review.ID = r.nextID
	review.CreatedAt = now
	review.UpdatedAt = now
	r.nextID++
	r.reviews = append(r.reviews, review)
	return review
}

func (r *productReviewRepo) FindByProductID(productID string) []model.ProductReview {
	r.delay()
	return r.filter(func(review model.ProductReview) bool { return review.ProductID == productID })
}

func (r *productReviewRepo) FindTop5ByProductID(productID string) []model.ProductReview {
	reviews := r.FindByProductIDOrderByCreatedAtDesc(productID)
	if len(reviews) > 5 {
		return reviews[:5]
	}
	return reviews
}

func (r *productReviewRepo) FindByProductIDs(productIDs []string) []model.ProductReview {
	r.delay()
	lookup := make(map[string]struct{}, len(productIDs))
	for _, productID := range productIDs {
		lookup[productID] = struct{}{}
	}
	return r.filter(func(review model.ProductReview) bool {
		_, ok := lookup[review.ProductID]
		return ok
	})
}

func (r *productReviewRepo) FindByCustomerID(customerID string) []model.ProductReview {
	r.delay()
	reviews := r.filter(func(review model.ProductReview) bool { return review.CustomerID == customerID })
	sortDesc(reviews)
	return reviews
}

func (r *productReviewRepo) FindByProductIDOrderByCreatedAtDesc(productID string) []model.ProductReview {
	r.delay()
	reviews := r.filter(func(review model.ProductReview) bool { return review.ProductID == productID })
	sortDesc(reviews)
	return reviews
}

func (r *productReviewRepo) FindRecent(limit int) []model.ProductReview {
	r.delay()
	reviews := r.filter(func(review model.ProductReview) bool { return true })
	sortDesc(reviews)
	if len(reviews) > limit {
		return reviews[:limit]
	}
	return reviews
}

func (r *productReviewRepo) FindByProductIDAndCreatedAtBetween(productID string, start time.Time, end time.Time) []model.ProductReview {
	r.delay()
	return r.filter(func(review model.ProductReview) bool {
		return review.ProductID == productID && !review.CreatedAt.Before(start) && review.CreatedAt.Before(end)
	})
}

func (r *productReviewRepo) FindByProductIDAndCreatedAtBetweenDesc(productID string, start time.Time, end time.Time) []model.ProductReview {
	reviews := r.FindByProductIDAndCreatedAtBetween(productID, start, end)
	sortDesc(reviews)
	return reviews
}

func (r *productReviewRepo) filter(predicate func(model.ProductReview) bool) []model.ProductReview {
	r.mu.RLock()
	defer r.mu.RUnlock()
	results := make([]model.ProductReview, 0)
	for _, review := range r.reviews {
		if predicate(review) {
			results = append(results, review)
		}
	}
	return results
}

func sortDesc(reviews []model.ProductReview) {
	sort.Slice(reviews, func(i, j int) bool {
		return reviews[i].CreatedAt.After(reviews[j].CreatedAt)
	})
}

func itoa(value int) string {
	return strconv.Itoa(value)
}
