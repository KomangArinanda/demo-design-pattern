package com.example.review.command.productreview;

import com.example.review.dto.command.request.GetListCustomerReviewsCommandRequest;
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

class GetListCustomerReviewsCommandTest {

    private final ProductReviewRepository repository = Mockito.mock(ProductReviewRepository.class);

    private final GetListCustomerReviewsCommand command =
        new GetListCustomerReviewsCommand(new DatabaseLatencySimulator(0), repository);

    @Test
    void executesCustomerReviewLookup() {
        when(repository.findByCustomerIdOrderByCreatedAtDesc("CUST-001"))
            .thenReturn(List.of(review(1L, "PROD-001", "CUST-001", "ORD-001", 5, "Great", at(2026, 5, 15, 10))));

        command.validate(new GetListCustomerReviewsCommandRequest("CUST-001"));
        List<ReviewDetailResponse> response = command.execute(new GetListCustomerReviewsCommandRequest("CUST-001"));

        assertEquals(1, response.size());
        assertEquals("CUST-001", response.get(0).customerId());
    }
}
