package com.example.review.controller;

import com.example.review.dto.request.CreateReviewRequest;
import com.example.review.dto.response.*;
import com.example.review.service.ProductReviewService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;
import javax.validation.constraints.Max;
import javax.validation.constraints.Min;
import javax.validation.constraints.NotBlank;
import java.util.List;

@Validated
@RestController
@RequestMapping("/api/v1")
@RequiredArgsConstructor
public class ProductReviewController {

    private final ProductReviewService service;

    @PostMapping("/products/{productId}/reviews")
    public ResponseEntity<CreateReviewResponse> createReview(@PathVariable @NotBlank String productId,
        @Valid @RequestBody CreateReviewRequest request) {
        return ResponseEntity.status(HttpStatus.CREATED).body(service.createReview(productId, request));
    }

    @GetMapping("/customers/{customerId}/reviews")
    public List<ReviewDetailResponse> getCustomerReviews(@PathVariable @NotBlank String customerId) {
        return service.getCustomerReviews(customerId);
    }

    @GetMapping("/products/{productId}/review-analytics/daily")
    public DailyProductReviewAnalyticsResponse getDailyProductReviewAnalytics(@PathVariable @NotBlank String productId,
        @RequestParam @Min(1) @Max(12) int month, @RequestParam @Min(1) int year) {
        return service.getDailyProductReviewAnalytics(productId, month, year);
    }

    @GetMapping("/products/{productId}/review-summary")
    public ProductReviewSummaryResponse getProductReviewSummary(@PathVariable @NotBlank String productId) {
        return service.getProductReviewSummary(productId);
    }

    @GetMapping("/products/{productId}/reviews")
    public List<ReviewDetailResponse> getProductReviews(@PathVariable @NotBlank String productId) {
        return service.getProductReviews(productId);
    }

    @GetMapping("/reviews/recent")
    public List<ReviewDetailResponse> getRecentReviews(@RequestParam(defaultValue = "10") int limit) {
        return service.getRecentReviews(limit);
    }

    @GetMapping("/sellers/{sellerId}/review-analytics")
    public SellerReviewAnalyticsResponse getSellerReviewAnalytics(@PathVariable @NotBlank String sellerId) {
        return service.getSellerReviewAnalytics(sellerId);
    }
}
