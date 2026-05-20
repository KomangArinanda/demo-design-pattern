package com.example.review.command.productreview;

import com.example.review.client.OrderClient;
import com.example.review.command.Command;
import com.example.review.dto.command.request.SaveProductReviewCommandRequest;
import com.example.review.config.DatabaseLatencySimulator;
import com.example.review.dto.response.CreateReviewResponse;
import com.example.review.dto.response.OrderValidationResponse;
import com.example.review.entity.ProductReview;
import com.example.review.exception.DuplicateReviewException;
import com.example.review.exception.InvalidOrderException;
import com.example.review.repository.ProductReviewRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Transactional;

@Component
@RequiredArgsConstructor
public class SaveProductReviewCommand implements Command<SaveProductReviewCommandRequest, CreateReviewResponse> {

    private final DatabaseLatencySimulator databaseLatencySimulator;

    private final OrderClient orderClient;

    private final ProductReviewRepository repository;

    @Override
    @Transactional
    public CreateReviewResponse execute(SaveProductReviewCommandRequest request) {
        databaseLatencySimulator.delay();
        ProductReview savedReview = repository.save(new ProductReview(
            request.productId(),
            request.request().customerId(),
            request.request().orderId(),
            request.request().rating(),
            request.request().comment()
        ));

        return new CreateReviewResponse(
            savedReview.getId(),
            savedReview.getProductId(),
            savedReview.getCustomerId(),
            savedReview.getOrderId(),
            savedReview.getRating(),
            savedReview.getComment(),
            savedReview.getCreatedAt()
        );
    }

    @Override
    public void validate(SaveProductReviewCommandRequest request) {
        OrderValidationResponse orderValidation = orderClient.validateOrder(
            request.request().customerId(),
            request.request().orderId(),
            request.productId()
        );

        if (!orderValidation.valid() || !"COMPLETED".equals(orderValidation.orderStatus())) {
            throw new InvalidOrderException();
        }

        databaseLatencySimulator.delay();
        if (repository.existsByProductIdAndCustomerIdAndOrderId(
            request.productId(),
            request.request().customerId(),
            request.request().orderId()
        )) {
            throw new DuplicateReviewException();
        }
    }
}
