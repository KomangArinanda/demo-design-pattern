package com.example.review.dto.response;

import java.time.LocalDateTime;

public record LatestReviewResponse(
    Long reviewId,
    String customerId,
    Integer rating,
    String comment,
    LocalDateTime createdAt
) {

}
