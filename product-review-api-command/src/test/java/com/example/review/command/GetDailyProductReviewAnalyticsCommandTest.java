package com.example.review.command;

import com.example.review.command.request.GetDailyProductReviewAnalyticsCommandRequest;
import com.example.review.config.DatabaseLatencySimulator;
import com.example.review.dto.response.DailyProductReviewAnalyticsResponse;
import com.example.review.repository.ProductReviewRepository;
import com.example.review.service.ProductReviewService;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.util.List;

import static com.example.review.command.support.ProductReviewTestFactory.at;
import static com.example.review.command.support.ProductReviewTestFactory.review;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.Mockito.when;

class GetDailyProductReviewAnalyticsCommandTest {

    private final ProductReviewRepository repository = Mockito.mock(ProductReviewRepository.class);

    private final ProductReviewService service = new ProductReviewService();

    private final GetDailyProductReviewAnalyticsCommand command =
        new GetDailyProductReviewAnalyticsCommand(new DatabaseLatencySimulator(0), repository, service);

    @Test
    void executesDailyAnalytics() {
        GetDailyProductReviewAnalyticsCommandRequest request =
            new GetDailyProductReviewAnalyticsCommandRequest("PROD-001", 5, 2026);
        when(repository.findByProductIdAndCreatedAtGreaterThanEqualAndCreatedAtLessThan(
            "PROD-001",
            at(2026, 5, 1, 0),
            at(2026, 6, 1, 0)
        )).thenReturn(List.of(
            review(1L, "PROD-001", "CUST-001", "ORD-001", 5, "Great", at(2026, 5, 1, 9)),
            review(2L, "PROD-001", "CUST-002", "ORD-002", 3, "Okay", at(2026, 5, 1, 12)),
            review(3L, "PROD-001", "CUST-003", "ORD-003", 1, "Bad", at(2026, 5, 2, 8))
        ));

        command.validate(request);
        DailyProductReviewAnalyticsResponse response = command.execute(request);

        assertEquals(3, response.totalReviews());
        assertEquals(3.0, response.averageRating());
        assertEquals(2, response.dailySummaries().size());
    }
}
