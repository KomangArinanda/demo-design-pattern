package com.example.review.config;

import com.example.review.entity.ProductReview;
import com.example.review.repository.ProductReviewRepository;
import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.CommandLineRunner;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.ArrayList;
import java.util.List;

@Configuration
@Slf4j
public class SeedDataConfig {

    private void seedProduct(List<ProductReview> reviews, String productId, int count, List<Integer> ratings) {
        for (int i = 1; i <= count; i++) {
            reviews.add(new ProductReview(
                productId,
                "CUST-" + productId.substring(productId.length() - 3) + "-" + i,
                "ORD-" + productId.substring(productId.length() - 3) + "-" + i,
                ratings.get(i - 1),
                "Sample review " + i + " for " + productId
            ));
        }
    }

    @Bean
    CommandLineRunner seedReviews(ProductReviewRepository repository) {
        return args -> {
            if (repository.count() > 0) {
                return;
            }

            List<ProductReview> reviews = new ArrayList<>();
            seedProduct(reviews, "PROD-001", 10, List.of(5, 5, 4, 4, 5, 3, 4, 5, 2, 5));
            seedProduct(reviews, "PROD-002", 8, List.of(5, 4, 4, 3, 2, 2, 1, 4));
            seedProduct(reviews, "PROD-003", 5, List.of(5, 5, 4, 2, 1));
            seedProduct(reviews, "PROD-004", 7, List.of(5, 4, 4, 4, 3, 2, 5));
            seedProduct(reviews, "PROD-005", 4, List.of(5, 4, 3, 1));

            repository.saveAll(reviews);
            log.info("Seeded {} product reviews", reviews.size());
        };
    }
}
