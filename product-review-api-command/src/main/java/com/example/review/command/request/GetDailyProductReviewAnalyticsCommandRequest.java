package com.example.review.command.request;

public record GetDailyProductReviewAnalyticsCommandRequest(
    String productId,
    int month,
    int year
) {

}
