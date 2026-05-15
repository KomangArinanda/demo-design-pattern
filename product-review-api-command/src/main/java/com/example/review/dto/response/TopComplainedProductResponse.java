package com.example.review.dto.response;

public record TopComplainedProductResponse(
    String productId,
    long negativeReviewCount
) {

}
