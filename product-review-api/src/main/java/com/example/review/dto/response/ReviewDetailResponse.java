package com.example.review.dto.response;

import java.time.LocalDateTime;

public record ReviewDetailResponse(
    Long reviewId,
    String productId,
    String customerId,
    String orderId,
    Integer rating,
    String comment,
    LocalDateTime createdAt
) {
}
