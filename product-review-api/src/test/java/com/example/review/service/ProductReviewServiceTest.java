package com.example.review.service;

import com.example.review.client.OrderClient;
import com.example.review.client.ProductClient;
import com.example.review.config.DatabaseLatencySimulator;
import com.example.review.dto.request.CreateReviewRequest;
import com.example.review.dto.response.*;
import com.example.review.entity.ProductReview;
import com.example.review.exception.DuplicateReviewException;
import com.example.review.exception.InvalidOrderException;
import com.example.review.repository.ProductReviewRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.data.domain.PageRequest;
import org.springframework.test.util.ReflectionTestUtils;

import java.time.LocalDateTime;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class ProductReviewServiceTest {

    @Mock
    private OrderClient orderClient;

    @Mock
    private ProductClient productClient;

    @Mock
    private ProductReviewRepository repository;

    private ProductReviewService service;

    private LocalDateTime at(int year, int month, int day, int hour) {
        return LocalDateTime.of(year, month, day, hour, 0);
    }

    @Test
    void createReviewCreatesReviewWhenOrderIsValidAndUnique() {
        CreateReviewRequest request = new CreateReviewRequest("CUST-001", "ORD-001", 5, "Great product");
        ProductReview savedReview = review(1L, "PROD-001", "CUST-001", "ORD-001", 5, "Great product", at(2026, 5, 15, 10));

        when(orderClient.validateOrder("CUST-001", "ORD-001", "PROD-001"))
            .thenReturn(new OrderValidationResponse(true, "COMPLETED"));
        when(repository.existsByProductIdAndCustomerIdAndOrderId("PROD-001", "CUST-001", "ORD-001"))
            .thenReturn(false);
        when(repository.save(any(ProductReview.class))).thenReturn(savedReview);

        CreateReviewResponse response = service.createReview("PROD-001", request);

        assertEquals(1L, response.reviewId());
        assertEquals("PROD-001", response.productId());
        assertEquals(5, response.rating());
        verify(repository).save(any(ProductReview.class));
    }

    @Test
    void createReviewRejectsDuplicateReview() {
        CreateReviewRequest request = new CreateReviewRequest("CUST-001", "ORD-001", 4, "Good");
        when(orderClient.validateOrder("CUST-001", "ORD-001", "PROD-001"))
            .thenReturn(new OrderValidationResponse(true, "COMPLETED"));
        when(repository.existsByProductIdAndCustomerIdAndOrderId("PROD-001", "CUST-001", "ORD-001"))
            .thenReturn(true);

        assertThrows(DuplicateReviewException.class, () -> service.createReview("PROD-001", request));
    }

    @Test
    void createReviewRejectsInvalidOrder() {
        CreateReviewRequest request = new CreateReviewRequest("CUST-001", "ORD-INVALID", 4, "Average");
        when(orderClient.validateOrder("CUST-001", "ORD-INVALID", "PROD-001"))
            .thenReturn(new OrderValidationResponse(false, "NOT_FOUND"));

        assertThrows(InvalidOrderException.class, () -> service.createReview("PROD-001", request));
    }

    @Test
    void getCustomerReviewsMapsRepositoryResults() {
        when(repository.findByCustomerIdOrderByCreatedAtDesc("CUST-001"))
            .thenReturn(List.of(review(1L, "PROD-001", "CUST-001", "ORD-001", 5, "Great", at(2026, 5, 15, 10))));

        List<ReviewDetailResponse> reviews = service.getCustomerReviews("CUST-001");

        assertEquals(1, reviews.size());
        assertEquals("PROD-001", reviews.get(0).productId());
        assertEquals("CUST-001", reviews.get(0).customerId());
    }

    @Test
    void getDailyProductReviewAnalyticsGroupsAndSummarizesByDay() {
        ProductReview first = review(1L, "PROD-001", "CUST-001", "ORD-001", 5, "Great", at(2026, 5, 1, 9));
        ProductReview second = review(2L, "PROD-001", "CUST-002", "ORD-002", 3, "Okay", at(2026, 5, 1, 12));
        ProductReview third = review(3L, "PROD-001", "CUST-003", "ORD-003", 1, "Bad", at(2026, 5, 2, 8));

        when(repository.findByProductIdAndCreatedAtGreaterThanEqualAndCreatedAtLessThan(
            "PROD-001",
            at(2026, 5, 1, 0),
            at(2026, 6, 1, 0)
        )).thenReturn(List.of(first, second, third));
        when(repository.findByProductIdAndCreatedAtGreaterThanEqualAndCreatedAtLessThanOrderByCreatedAtDesc(
            "PROD-001",
            at(2026, 5, 1, 0),
            at(2026, 6, 1, 0)
        )).thenReturn(List.of(third, second, first));

        DailyProductReviewAnalyticsResponse response = service.getDailyProductReviewAnalytics("PROD-001", 5, 2026);

        assertEquals(3, response.totalReviews());
        assertEquals(3.0, response.averageRating());
        assertEquals(2, response.dailySummaries().size());
        assertEquals(2, response.dailySummaries().get(0).totalReviews());
        assertEquals(1, response.dailySummaries().get(1).totalReviews());
    }

    @Test
    void getProductReviewSummaryCalculatesSummary() {
        List<ProductReview> reviews = List.of(
            review(1L, "PROD-001", "CUST-001", "ORD-001", 5, "Great", at(2026, 5, 15, 10)),
            review(2L, "PROD-001", "CUST-002", "ORD-002", 3, "Okay", at(2026, 5, 15, 9))
        );
        when(repository.findByProductId("PROD-001")).thenReturn(reviews);
        when(repository.findTop5ByProductIdOrderByCreatedAtDesc("PROD-001")).thenReturn(reviews);

        ProductReviewSummaryResponse response = service.getProductReviewSummary("PROD-001");

        assertEquals(2, response.totalReviews());
        assertEquals(4.0, response.averageRating());
        assertEquals(1, response.ratingDistribution().get(5));
        assertEquals(2, response.latestReviews().size());
    }

    @Test
    void getProductReviewsMapsRepositoryResults() {
        when(repository.findByProductIdOrderByCreatedAtDesc("PROD-001"))
            .thenReturn(List.of(review(1L, "PROD-001", "CUST-001", "ORD-001", 5, "Great", at(2026, 5, 15, 10))));

        List<ReviewDetailResponse> reviews = service.getProductReviews("PROD-001");

        assertEquals(1, reviews.size());
        assertEquals("PROD-001", reviews.get(0).productId());
    }

    @Test
    void getRecentReviewsNormalizesLimitAndMapsResults() {
        when(repository.findAllByOrderByCreatedAtDesc(PageRequest.of(0, 1)))
            .thenReturn(List.of(review(1L, "PROD-001", "CUST-001", "ORD-001", 5, "Great", at(2026, 5, 15, 10))));

        List<ReviewDetailResponse> reviews = service.getRecentReviews(0);

        assertEquals(1, reviews.size());
        verify(repository).findAllByOrderByCreatedAtDesc(PageRequest.of(0, 1));
    }

    @Test
    void getSellerReviewAnalyticsCalculatesSellerMetrics() {
        when(productClient.getProductIdsBySellerId("SELLER-001"))
            .thenReturn(new SellerProductsResponse("SELLER-001", List.of("PROD-001", "PROD-002")));
        when(repository.findByProductIdIn(List.of("PROD-001", "PROD-002")))
            .thenReturn(List.of(
                review(1L, "PROD-001", "CUST-001", "ORD-001", 5, "Great", at(2026, 5, 15, 10)),
                review(2L, "PROD-002", "CUST-002", "ORD-002", 2, "Bad", at(2026, 5, 15, 9)),
                review(3L, "PROD-002", "CUST-003", "ORD-003", 1, "Worse", at(2026, 5, 15, 8))
            ));

        SellerReviewAnalyticsResponse response = service.getSellerReviewAnalytics("SELLER-001");

        assertEquals(3, response.totalReviews());
        assertEquals(2.7, response.averageRating());
        assertEquals(2, response.negativeReviewCount());
        assertEquals(66.7, response.negativeReviewPercentage());
        assertEquals("PROD-002", response.topComplainedProducts().get(0).productId());
    }

    private ProductReview review(
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

    @BeforeEach
    void setUp() {
        service = new ProductReviewService(new DatabaseLatencySimulator(0), orderClient, productClient, repository);
    }
}
