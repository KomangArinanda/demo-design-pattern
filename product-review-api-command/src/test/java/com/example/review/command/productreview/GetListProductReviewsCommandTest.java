package com.example.review.command.productreview;

import com.example.review.dto.command.request.GetListProductReviewsCommandRequest;
import com.example.review.config.DatabaseLatencySimulator;
import com.example.review.dto.response.ReviewDetailResponse;
import com.example.review.repository.ProductReviewRepository;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.util.List;

import static com.example.review.command.support.ProductReviewTestFactory.at;
import static com.example.review.command.support.ProductReviewTestFactory.review;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.Mockito.when;

class GetListProductReviewsCommandTest {

    private final ProductReviewRepository repository = Mockito.mock(ProductReviewRepository.class);

    private final GetListProductReviewsCommand command =
        new GetListProductReviewsCommand(new DatabaseLatencySimulator(0), repository);

    @Test
    void executesProductReviewLookup() {
        when(repository.findByProductIdOrderByCreatedAtDesc("PROD-001"))
            .thenReturn(List.of(review(1L, "PROD-001", "CUST-001", "ORD-001", 5, "Great", at(2026, 5, 15, 10))));

        command.validate(new GetListProductReviewsCommandRequest("PROD-001"));
        List<ReviewDetailResponse> response = command.execute(new GetListProductReviewsCommandRequest("PROD-001"));

        assertEquals(1, response.size());
        assertEquals("PROD-001", response.get(0).productId());
    }
}
