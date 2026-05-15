package com.example.review.dto.request;

import javax.validation.constraints.*;

public record CreateReviewRequest(
    @NotBlank(message = "Customer ID is required")
    String customerId,
    @NotBlank(message = "Order ID is required")
    String orderId,
    @NotNull(message = "Rating is required")
    @Min(value = 1, message = "Rating must be between 1 and 5")
    @Max(value = 5, message = "Rating must be between 1 and 5")
    Integer rating,
    @Size(max = 500, message = "Comment must be at most 500 characters")
    String comment
) {

}
