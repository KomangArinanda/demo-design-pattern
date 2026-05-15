package product_review

import (
	"context"
	"example/product-review-api-command-go/internal/client"
	"example/product-review-api-command-go/internal/dto/response"
	"example/product-review-api-command-go/internal/repo"
	"example/product-review-api-command-go/internal/service"
	"example/product-review-api-command-go/internal/shared/appctx"
	"example/product-review-api-command-go/internal/usecase/common"
	"sort"
)

type getSellerReviewAnalyticsUsecase struct {
	repository    repo.ProductReviewRepo
	productClient client.ProductClient
	service       *service.ProductReviewService
}

func NewGetSellerReviewAnalyticsUsecase(repository repo.ProductReviewRepo, productClient client.ProductClient, sharedService *service.ProductReviewService) *getSellerReviewAnalyticsUsecase {
	return &getSellerReviewAnalyticsUsecase{
		repository:    repository,
		productClient: productClient,
		service:       sharedService,
	}
}

func (u *getSellerReviewAnalyticsUsecase) Execute(_ context.Context, input any) appctx.Response {
	request, ok := common.MustInput[GetSellerReviewAnalyticsUsecaseRequest](input)
	if !ok {
		return common.BadRequest("Invalid request")
	}

	sellerProducts := u.productClient.GetProductIDsBySellerID(request.SellerID)
	reviews := u.repository.FindByProductIDs(sellerProducts.ProductIDs)
	negativeCount := 0
	negativeByProduct := map[string]int{}
	for _, review := range reviews {
		if review.Rating <= 2 {
			negativeCount++
			negativeByProduct[review.ProductID]++
		}
	}
	topComplainedProducts := make([]response.TopComplainedProductResponse, 0, len(negativeByProduct))
	for productID, count := range negativeByProduct {
		topComplainedProducts = append(topComplainedProducts, response.TopComplainedProductResponse{
			ProductID:           productID,
			NegativeReviewCount: count,
		})
	}
	sort.Slice(topComplainedProducts, func(i, j int) bool {
		if topComplainedProducts[i].NegativeReviewCount == topComplainedProducts[j].NegativeReviewCount {
			return topComplainedProducts[i].ProductID < topComplainedProducts[j].ProductID
		}
		return topComplainedProducts[i].NegativeReviewCount > topComplainedProducts[j].NegativeReviewCount
	})
	if len(topComplainedProducts) > 3 {
		topComplainedProducts = topComplainedProducts[:3]
	}
	return appctx.OK(response.SellerReviewAnalyticsResponse{
		SellerID:                 request.SellerID,
		TotalReviews:             len(reviews),
		AverageRating:            u.service.CalculateAverageRating(reviews),
		NegativeReviewCount:      negativeCount,
		NegativeReviewPercentage: u.service.CalculatePercentage(negativeCount, len(reviews)),
		RatingDistribution:       u.service.CalculateRatingDistribution(reviews),
		TopComplainedProducts:    topComplainedProducts,
	})
}
