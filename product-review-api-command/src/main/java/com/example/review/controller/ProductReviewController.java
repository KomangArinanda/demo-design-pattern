package com.example.review.controller;

import com.example.review.command.*;
import com.example.review.command.request.*;
import com.example.review.dto.request.CreateReviewRequest;
import com.example.review.dto.response.*;
import jakarta.validation.Valid;
import jakarta.validation.constraints.Max;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@Validated
@RestController
@RequestMapping("/api/v1")
@RequiredArgsConstructor
public class ProductReviewController {

    private final CommandExecutor commandExecutor;

    @PostMapping("/products/{productId}/reviews")
    public ResponseEntity<CreateReviewResponse> createReview(@PathVariable @NotBlank String productId,
        @Valid @RequestBody CreateReviewRequest request) {
        return ResponseEntity.status(HttpStatus.CREATED).body(commandExecutor.execute(SaveProductReviewCommand.class, new SaveProductReviewCommandRequest(productId, request)));
    }

    @GetMapping("/customers/{customerId}/reviews")
    public List<ReviewDetailResponse> getCustomerReviews(@PathVariable @NotBlank String customerId) {
        return commandExecutor.execute(GetListCustomerReviewsCommand.class, new GetListCustomerReviewsCommandRequest(customerId));
    }

    @GetMapping("/products/{productId}/review-analytics/daily")
    public DailyProductReviewAnalyticsResponse getDailyProductReviewAnalytics(@PathVariable @NotBlank String productId,
        @RequestParam @Min(1) @Max(12) int month, @RequestParam @Min(1) int year) {
        return commandExecutor.execute(GetDailyProductReviewAnalyticsCommand.class, new GetDailyProductReviewAnalyticsCommandRequest(productId, month, year));
    }

    @GetMapping("/products/{productId}/review-summary")
    public ProductReviewSummaryResponse getProductReviewSummary(@PathVariable @NotBlank String productId) {
        return commandExecutor.execute(GetProductReviewSummaryCommand.class, new GetProductReviewSummaryCommandRequest(productId));
    }

    @GetMapping("/products/{productId}/reviews")
    public List<ReviewDetailResponse> getProductReviews(@PathVariable @NotBlank String productId) {
        return commandExecutor.execute(GetListProductReviewsCommand.class, new GetListProductReviewsCommandRequest(productId));
    }

    @GetMapping("/reviews/recent")
    public List<ReviewDetailResponse> getRecentReviews(@RequestParam(defaultValue = "10") int limit) {
        return commandExecutor.execute(GetListRecentReviewsCommand.class, new GetListRecentReviewsCommandRequest(limit));
    }

    @GetMapping("/sellers/{sellerId}/review-analytics")
    public SellerReviewAnalyticsResponse getSellerReviewAnalytics(@PathVariable @NotBlank String sellerId) {
        return commandExecutor.execute(GetSellerReviewAnalyticsCommand.class, new GetSellerReviewAnalyticsCommandRequest(sellerId));
    }
}
