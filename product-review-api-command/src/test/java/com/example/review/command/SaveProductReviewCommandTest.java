package com.example.review.command;

import com.example.review.client.OrderClient;
import com.example.review.command.request.SaveProductReviewCommandRequest;
import com.example.review.config.DatabaseLatencySimulator;
import com.example.review.dto.request.CreateReviewRequest;
import com.example.review.dto.response.CreateReviewResponse;
import com.example.review.dto.response.OrderValidationResponse;
import com.example.review.exception.DuplicateReviewException;
import com.example.review.exception.InvalidOrderException;
import com.example.review.repository.ProductReviewRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import static com.example.review.command.support.ProductReviewTestFactory.at;
import static com.example.review.command.support.ProductReviewTestFactory.review;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class SaveProductReviewCommandTest {

    private SaveProductReviewCommand command;

    @Mock
    private OrderClient orderClient;

    @Mock
    private ProductReviewRepository repository;

    @Test
    void rejectsDuplicateReview() {
        SaveProductReviewCommandRequest request = new SaveProductReviewCommandRequest(
            "PROD-001",
            new CreateReviewRequest("CUST-001", "ORD-001", 4, "Good")
        );
        when(orderClient.validateOrder("CUST-001", "ORD-001", "PROD-001"))
            .thenReturn(new OrderValidationResponse(true, "COMPLETED"));
        when(repository.existsByProductIdAndCustomerIdAndOrderId("PROD-001", "CUST-001", "ORD-001"))
            .thenReturn(true);

        assertThrows(DuplicateReviewException.class, () -> command.validate(request));
    }

    @Test
    void rejectsInvalidOrder() {
        SaveProductReviewCommandRequest request = new SaveProductReviewCommandRequest(
            "PROD-001",
            new CreateReviewRequest("CUST-001", "ORD-INVALID", 4, "Average")
        );
        when(orderClient.validateOrder("CUST-001", "ORD-INVALID", "PROD-001"))
            .thenReturn(new OrderValidationResponse(false, "NOT_FOUND"));

        assertThrows(InvalidOrderException.class, () -> command.validate(request));
    }

    @BeforeEach
    void setUp() {
        command = new SaveProductReviewCommand(new DatabaseLatencySimulator(0), orderClient, repository);
    }

    @Test
    void validatesAndExecutesSave() {
        SaveProductReviewCommandRequest request = new SaveProductReviewCommandRequest(
            "PROD-001",
            new CreateReviewRequest("CUST-001", "ORD-001", 5, "Great")
        );
        when(orderClient.validateOrder("CUST-001", "ORD-001", "PROD-001"))
            .thenReturn(new OrderValidationResponse(true, "COMPLETED"));
        when(repository.existsByProductIdAndCustomerIdAndOrderId("PROD-001", "CUST-001", "ORD-001"))
            .thenReturn(false);
        when(repository.save(any())).thenReturn(review(1L, "PROD-001", "CUST-001", "ORD-001", 5, "Great", at(2026, 5, 15, 10)));

        command.validate(request);
        CreateReviewResponse response = command.execute(request);

        assertEquals(1L, response.reviewId());
        assertEquals("PROD-001", response.productId());
    }
}
