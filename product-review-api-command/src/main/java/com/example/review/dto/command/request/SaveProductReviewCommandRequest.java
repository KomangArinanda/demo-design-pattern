package com.example.review.dto.command.request;

import com.example.review.dto.request.CreateReviewRequest;

public record SaveProductReviewCommandRequest(
    String productId,
    CreateReviewRequest request
) {

}
