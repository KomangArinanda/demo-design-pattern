package com.example.review.dto.command.request;

public record GetDailyProductReviewAnalyticsCommandRequest(
    String productId,
    int month,
    int year
) {

}
