package com.example.review.command;

import com.example.review.command.request.GetListRecentReviewsCommandRequest;
import com.example.review.config.DatabaseLatencySimulator;
import com.example.review.dto.response.ReviewDetailResponse;
import com.example.review.repository.ProductReviewRepository;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import org.springframework.data.domain.PageRequest;

import java.util.List;

import static com.example.review.command.support.ProductReviewTestFactory.at;
import static com.example.review.command.support.ProductReviewTestFactory.review;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

class GetListRecentReviewsCommandTest {

    private final ProductReviewRepository repository = Mockito.mock(ProductReviewRepository.class);

    private final GetListRecentReviewsCommand command =
        new GetListRecentReviewsCommand(new DatabaseLatencySimulator(0), repository);

    @Test
    void executesRecentReviewLookupWithNormalizedLimit() {
        when(repository.findAllByOrderByCreatedAtDesc(PageRequest.of(0, 1)))
            .thenReturn(List.of(review(1L, "PROD-001", "CUST-001", "ORD-001", 5, "Great", at(2026, 5, 15, 10))));

        command.validate(new GetListRecentReviewsCommandRequest(0));
        List<ReviewDetailResponse> response = command.execute(new GetListRecentReviewsCommandRequest(0));

        assertEquals(1, response.size());
        verify(repository).findAllByOrderByCreatedAtDesc(PageRequest.of(0, 1));
    }
}
