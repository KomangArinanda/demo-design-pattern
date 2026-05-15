package com.example.review.dto.response;

public record OrderValidationResponse(
    boolean valid,
    String orderStatus
) {

}
