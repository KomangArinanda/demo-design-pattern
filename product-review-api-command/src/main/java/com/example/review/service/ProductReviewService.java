package com.example.review.service;

import com.example.review.entity.ProductReview;
import org.springframework.stereotype.Service;

import java.math.BigDecimal;
import java.math.RoundingMode;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

@Service
public class ProductReviewService {

    public double calculateAverageRating(List<ProductReview> reviews) {
        if (reviews.isEmpty()) {
            return 0.0;
        }
        double average = reviews.stream()
            .mapToInt(ProductReview::getRating)
            .average()
            .orElse(0.0);
        return round(average);
    }

    public Map<Integer, Long> calculateRatingDistribution(List<ProductReview> reviews) {
        Map<Integer, Long> distribution = new LinkedHashMap<>();
        for (int rating = 1; rating <= 5; rating++) {
            distribution.put(rating, 0L);
        }
        reviews.forEach(review -> distribution.computeIfPresent(review.getRating(), (rating, count) -> count + 1));
        return distribution;
    }

    private double round(double value) {
        return BigDecimal.valueOf(value)
            .setScale(1, RoundingMode.HALF_UP)
            .doubleValue();
    }

}
