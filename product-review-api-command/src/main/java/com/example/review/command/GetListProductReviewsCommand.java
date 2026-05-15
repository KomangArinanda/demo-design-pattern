package com.example.review.command;

import com.example.review.command.request.GetListProductReviewsCommandRequest;
import com.example.review.config.DatabaseLatencySimulator;
import com.example.review.dto.response.ReviewDetailResponse;
import com.example.review.entity.ProductReview;
import com.example.review.repository.ProductReviewRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;

import java.util.List;

@Component
@RequiredArgsConstructor
public class GetListProductReviewsCommand implements Command<GetListProductReviewsCommandRequest, List<ReviewDetailResponse>> {

    private final DatabaseLatencySimulator databaseLatencySimulator;

    private final ProductReviewRepository repository;

    @Override
    public List<ReviewDetailResponse> execute(GetListProductReviewsCommandRequest request) {
        databaseLatencySimulator.delay();
        return repository.findByProductIdOrderByCreatedAtDesc(request.productId())
            .stream()
            .map(this::toReviewDetailResponse)
            .toList();
    }

    private ReviewDetailResponse toReviewDetailResponse(ProductReview review) {
        return new ReviewDetailResponse(
            review.getId(),
            review.getProductId(),
            review.getCustomerId(),
            review.getOrderId(),
            review.getRating(),
            review.getComment(),
            review.getCreatedAt()
        );
    }

}
