package com.example.review.command.productreview;

import com.example.review.dto.command.request.GetProductReviewSummaryCommandRequest;
import com.example.review.config.DatabaseLatencySimulator;
import com.example.review.dto.response.ProductReviewSummaryResponse;
import com.example.review.entity.ProductReview;
import com.example.review.repository.ProductReviewRepository;
import com.example.review.service.ProductReviewService;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.util.List;

import static com.example.review.command.support.ProductReviewTestFactory.at;
import static com.example.review.command.support.ProductReviewTestFactory.review;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.Mockito.when;

class GetProductReviewSummaryCommandTest {

    private final ProductReviewRepository repository = Mockito.mock(ProductReviewRepository.class);

    private final ProductReviewService service = new ProductReviewService();

    private final GetProductReviewSummaryCommand command =
        new GetProductReviewSummaryCommand(new DatabaseLatencySimulator(0), repository, service);

    @Test
    void executesProductSummary() {
        List<ProductReview> reviews = List.of(
            review(1L, "PROD-001", "CUST-001", "ORD-001", 5, "Great", at(2026, 5, 15, 10)),
            review(2L, "PROD-001", "CUST-002", "ORD-002", 3, "Okay", at(2026, 5, 15, 9))
        );
        when(repository.findByProductId("PROD-001")).thenReturn(reviews);
        when(repository.findTop5ByProductIdOrderByCreatedAtDesc("PROD-001")).thenReturn(reviews);

        command.validate(new GetProductReviewSummaryCommandRequest("PROD-001"));
        ProductReviewSummaryResponse response = command.execute(new GetProductReviewSummaryCommandRequest("PROD-001"));

        assertEquals(2, response.totalReviews());
        assertEquals(4.0, response.averageRating());
    }
}
