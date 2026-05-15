package com.example.review.dto.response;

import java.util.List;
import java.util.Map;

public record DailyProductReviewAnalyticsResponse(
    String productId,
    int month,
    int year,
    long totalReviews,
    double averageRating,
    Map<Integer, Long> ratingDistribution,
    List<DailyReviewSummaryResponse> dailySummaries
) {

}
