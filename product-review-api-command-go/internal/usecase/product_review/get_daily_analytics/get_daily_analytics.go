package get_daily_analytics

import (
	"context"
	"errors"
	"example/product-review-api-command-go/internal/dto/response"
	"example/product-review-api-command-go/internal/model"
	"example/product-review-api-command-go/internal/repo"
	"example/product-review-api-command-go/internal/service"
	"example/product-review-api-command-go/internal/shared/appctx"
	"example/product-review-api-command-go/internal/usecase/common"
	"net/http"
	"sort"
	"time"
)

type GetDailyAnalyticsRequest struct {
	ProductID string
	Month     int
	Year      int
}

var ErrInvalidMonth = errors.New("month must be between 1 and 12")

func errorResponse(err error) appctx.Response {
	switch err {
	case ErrInvalidMonth:
		return appctx.Error(http.StatusBadRequest, err.Error())
	default:
		return appctx.Error(http.StatusBadRequest, err.Error())
	}
}

type getDailyAnalytics struct {
	repository repo.ProductReviewRepo
	service    *service.ProductReviewService
}

func NewGetDailyAnalytics(repository repo.ProductReviewRepo, sharedService *service.ProductReviewService) *getDailyAnalytics {
	return &getDailyAnalytics{
		repository: repository,
		service:    sharedService,
	}
}

func (u *getDailyAnalytics) Execute(_ context.Context, input any) appctx.Response {
	request, ok := common.MustInput[GetDailyAnalyticsRequest](input)
	if !ok {
		return common.BadRequest("Invalid request")
	}
	if request.Month < 1 || request.Month > 12 {
		return errorResponse(ErrInvalidMonth)
	}

	start := time.Date(request.Year, time.Month(request.Month), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0)
	monthlyReviews := u.repository.FindByProductIDAndCreatedAtBetween(request.ProductID, start, end)
	_ = u.repository.FindByProductIDAndCreatedAtBetweenDesc(request.ProductID, start, end)

	reviewsByDate := map[string][]model.ProductReview{}
	for _, review := range monthlyReviews {
		key := review.CreatedAt.Format("2006-01-02")
		reviewsByDate[key] = append(reviewsByDate[key], review)
	}
	keys := make([]string, 0, len(reviewsByDate))
	for key := range reviewsByDate {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	dailySummaries := make([]response.DailyReviewSummaryResponse, 0, len(keys))
	for _, key := range keys {
		reviews := reviewsByDate[key]
		dailySummaries = append(dailySummaries, response.DailyReviewSummaryResponse{
			Date:               key,
			TotalReviews:       len(reviews),
			AverageRating:      u.service.CalculateAverageRating(reviews),
			RatingDistribution: u.service.CalculateRatingDistribution(reviews),
		})
	}
	return appctx.OK(response.DailyProductReviewAnalyticsResponse{
		ProductID:          request.ProductID,
		Month:              request.Month,
		Year:               request.Year,
		TotalReviews:       len(monthlyReviews),
		AverageRating:      u.service.CalculateAverageRating(monthlyReviews),
		RatingDistribution: u.service.CalculateRatingDistribution(monthlyReviews),
		DailySummaries:     dailySummaries,
	})
}
