package com.example.review.dto.response;

import java.util.List;
import java.util.Map;

public record SellerReviewAnalyticsResponse(
    String sellerId,
    long totalReviews,
    double averageRating,
    long negativeReviewCount,
    double negativeReviewPercentage,
    Map<Integer, Long> ratingDistribution,
    List<TopComplainedProductResponse> topComplainedProducts
) {

}
