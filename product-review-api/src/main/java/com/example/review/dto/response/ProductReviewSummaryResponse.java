package com.example.review.dto.response;

import java.util.List;
import java.util.Map;

public record ProductReviewSummaryResponse(
    String productId,
    long totalReviews,
    double averageRating,
    Map<Integer, Long> ratingDistribution,
    List<LatestReviewResponse> latestReviews
) {

}
