package com.example.review.repository;

import com.example.review.entity.ProductReview;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.domain.Pageable;

import java.time.LocalDateTime;
import java.util.List;

public interface ProductReviewRepository extends JpaRepository<ProductReview, Long> {

    boolean existsByProductIdAndCustomerIdAndOrderId(String productId, String customerId, String orderId);

    List<ProductReview> findByProductId(String productId);

    List<ProductReview> findByCustomerIdOrderByCreatedAtDesc(String customerId);

    List<ProductReview> findByProductIdIn(List<String> productIds);

    List<ProductReview> findTop5ByProductIdOrderByCreatedAtDesc(String productId);

    List<ProductReview> findByProductIdOrderByCreatedAtDesc(String productId);

    List<ProductReview> findAllByOrderByCreatedAtDesc(Pageable pageable);

    List<ProductReview> findByProductIdAndCreatedAtGreaterThanEqualAndCreatedAtLessThan(
        String productId,
        LocalDateTime startDateTime,
        LocalDateTime endDateTime
    );

    List<ProductReview> findByProductIdAndCreatedAtGreaterThanEqualAndCreatedAtLessThanOrderByCreatedAtDesc(
        String productId,
        LocalDateTime startDateTime,
        LocalDateTime endDateTime
    );
}
