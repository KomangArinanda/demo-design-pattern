package service

import (
	"example/product-review-api-command-go/internal/dto/response"
	"example/product-review-api-command-go/internal/model"
	"math"
)

type ProductReviewService struct {
}

func NewProductReviewService() *ProductReviewService {
	return &ProductReviewService{}
}

func (s *ProductReviewService) CalculateAverageRating(reviews []model.ProductReview) float64 {
	if len(reviews) == 0 {
		return 0
	}
	total := 0
	for _, review := range reviews {
		total += review.Rating
	}
	return round(float64(total) / float64(len(reviews)))
}

func (s *ProductReviewService) CalculatePercentage(numerator int, denominator int) float64 {
	if denominator == 0 {
		return 0
	}
	return round(float64(numerator) * 100 / float64(denominator))
}

func (s *ProductReviewService) CalculateRatingDistribution(reviews []model.ProductReview) map[int]int {
	distribution := map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0}
	for _, review := range reviews {
		distribution[review.Rating]++
	}
	return distribution
}

func (s *ProductReviewService) NormalizeRecentReviewLimit(limit int) int {
	if limit < 1 {
		return 1
	}
	if limit > 50 {
		return 50
	}
	return limit
}

func (s *ProductReviewService) ToLatestReviewResponse(review model.ProductReview) response.LatestReviewResponse {
	return response.LatestReviewResponse{
		ReviewID:   review.ID,
		CustomerID: review.CustomerID,
		Rating:     review.Rating,
		Comment:    review.Comment,
		CreatedAt:  review.CreatedAt,
	}
}

func (s *ProductReviewService) MapLatestReviewResponses(reviews []model.ProductReview) []response.LatestReviewResponse {
	results := make([]response.LatestReviewResponse, 0, len(reviews))
	for _, review := range reviews {
		results = append(results, s.ToLatestReviewResponse(review))
	}
	return results
}

func (s *ProductReviewService) ToReviewDetailResponse(review model.ProductReview) response.ReviewDetailResponse {
	return response.ReviewDetailResponse{
		ReviewID:   review.ID,
		ProductID:  review.ProductID,
		CustomerID: review.CustomerID,
		OrderID:    review.OrderID,
		Rating:     review.Rating,
		Comment:    review.Comment,
		CreatedAt:  review.CreatedAt,
	}
}

func (s *ProductReviewService) MapReviewDetailResponses(reviews []model.ProductReview) []response.ReviewDetailResponse {
	results := make([]response.ReviewDetailResponse, 0, len(reviews))
	for _, review := range reviews {
		results = append(results, s.ToReviewDetailResponse(review))
	}
	return results
}

func round(value float64) float64 {
	return math.Round(value*10) / 10
}
