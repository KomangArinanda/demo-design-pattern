package com.example.review.dto.response;

import java.time.LocalDate;
import java.util.Map;

public record DailyReviewSummaryResponse(
    LocalDate date,
    long totalReviews,
    double averageRating,
    Map<Integer, Long> ratingDistribution
) {

}
