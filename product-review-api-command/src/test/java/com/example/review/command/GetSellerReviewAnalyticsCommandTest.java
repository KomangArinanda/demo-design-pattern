package com.example.review.command;

import com.example.review.client.ProductClient;
import com.example.review.command.request.GetSellerReviewAnalyticsCommandRequest;
import com.example.review.config.DatabaseLatencySimulator;
import com.example.review.dto.response.SellerProductsResponse;
import com.example.review.dto.response.SellerReviewAnalyticsResponse;
import com.example.review.repository.ProductReviewRepository;
import com.example.review.service.ProductReviewService;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.util.List;

import static com.example.review.command.support.ProductReviewTestFactory.at;
import static com.example.review.command.support.ProductReviewTestFactory.review;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.Mockito.when;

class GetSellerReviewAnalyticsCommandTest {

    private final ProductClient productClient = Mockito.mock(ProductClient.class);

    private final ProductReviewRepository repository = Mockito.mock(ProductReviewRepository.class);

    private final ProductReviewService service = new ProductReviewService();

    private final GetSellerReviewAnalyticsCommand command =
        new GetSellerReviewAnalyticsCommand(new DatabaseLatencySimulator(0), productClient, repository, service);

    @Test
    void executesSellerAnalytics() {
        when(productClient.getProductIdsBySellerId("SELLER-001"))
            .thenReturn(new SellerProductsResponse("SELLER-001", List.of("PROD-001", "PROD-002")));
        when(repository.findByProductIdIn(List.of("PROD-001", "PROD-002")))
            .thenReturn(List.of(
                review(1L, "PROD-001", "CUST-001", "ORD-001", 5, "Great", at(2026, 5, 15, 10)),
                review(2L, "PROD-002", "CUST-002", "ORD-002", 2, "Bad", at(2026, 5, 15, 9)),
                review(3L, "PROD-002", "CUST-003", "ORD-003", 1, "Worse", at(2026, 5, 15, 8))
            ));

        command.validate(new GetSellerReviewAnalyticsCommandRequest("SELLER-001"));
        SellerReviewAnalyticsResponse response =
            command.execute(new GetSellerReviewAnalyticsCommandRequest("SELLER-001"));

        assertEquals(3, response.totalReviews());
        assertEquals(2, response.negativeReviewCount());
        assertEquals("PROD-002", response.topComplainedProducts().get(0).productId());
    }
}
