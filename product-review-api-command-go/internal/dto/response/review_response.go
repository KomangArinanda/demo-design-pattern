package response

import "time"

type CreateReviewResponse struct {
	ReviewID   int64     `json:"reviewId"`
	ProductID  string    `json:"productId"`
	CustomerID string    `json:"customerId"`
	OrderID    string    `json:"orderId"`
	Rating     int       `json:"rating"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"createdAt"`
}

type LatestReviewResponse struct {
	ReviewID   int64     `json:"reviewId"`
	CustomerID string    `json:"customerId"`
	Rating     int       `json:"rating"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"createdAt"`
}

type ReviewDetailResponse struct {
	ReviewID   int64     `json:"reviewId"`
	ProductID  string    `json:"productId"`
	CustomerID string    `json:"customerId"`
	OrderID    string    `json:"orderId"`
	Rating     int       `json:"rating"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"createdAt"`
}

type ProductReviewSummaryResponse struct {
	ProductID          string                    `json:"productId"`
	TotalReviews       int                       `json:"totalReviews"`
	AverageRating      float64                   `json:"averageRating"`
	RatingDistribution map[int]int               `json:"ratingDistribution"`
	LatestReviews      []LatestReviewResponse    `json:"latestReviews"`
}

type SellerReviewAnalyticsResponse struct {
	SellerID                 string                         `json:"sellerId"`
	TotalReviews             int                            `json:"totalReviews"`
	AverageRating            float64                        `json:"averageRating"`
	NegativeReviewCount      int                            `json:"negativeReviewCount"`
	NegativeReviewPercentage float64                        `json:"negativeReviewPercentage"`
	RatingDistribution       map[int]int                    `json:"ratingDistribution"`
	TopComplainedProducts    []TopComplainedProductResponse `json:"topComplainedProducts"`
}

type TopComplainedProductResponse struct {
	ProductID           string `json:"productId"`
	NegativeReviewCount int    `json:"negativeReviewCount"`
}

type DailyReviewSummaryResponse struct {
	Date               string         `json:"date"`
	TotalReviews       int            `json:"totalReviews"`
	AverageRating      float64        `json:"averageRating"`
	RatingDistribution map[int]int    `json:"ratingDistribution"`
}

type DailyProductReviewAnalyticsResponse struct {
	ProductID          string                       `json:"productId"`
	Month              int                          `json:"month"`
	Year               int                          `json:"year"`
	TotalReviews       int                          `json:"totalReviews"`
	AverageRating      float64                      `json:"averageRating"`
	RatingDistribution map[int]int                  `json:"ratingDistribution"`
	DailySummaries     []DailyReviewSummaryResponse `json:"dailySummaries"`
}

type OrderValidationResponse struct {
	Valid       bool
	OrderStatus string
}

type SellerProductsResponse struct {
	SellerID   string
	ProductIDs []string
}
