package com.example.review.command.support;

import com.example.review.entity.ProductReview;
import org.springframework.test.util.ReflectionTestUtils;

import java.time.LocalDateTime;

public final class ProductReviewTestFactory {

    private ProductReviewTestFactory() {
    }

    public static ProductReview review(
        Long id,
        String productId,
        String customerId,
        String orderId,
        int rating,
        String comment,
        LocalDateTime createdAt
    ) {
        ProductReview review = new ProductReview(productId, customerId, orderId, rating, comment);
        ReflectionTestUtils.setField(review, "id", id);
        ReflectionTestUtils.setField(review, "createdAt", createdAt);
        ReflectionTestUtils.setField(review, "updatedAt", createdAt);
        return review;
    }

    public static LocalDateTime at(int year, int month, int day, int hour) {
        return LocalDateTime.of(year, month, day, hour, 0);
    }
}
