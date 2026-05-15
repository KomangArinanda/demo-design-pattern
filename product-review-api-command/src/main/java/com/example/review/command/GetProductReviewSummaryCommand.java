package com.example.review.command;

import com.example.review.command.request.GetProductReviewSummaryCommandRequest;
import com.example.review.config.DatabaseLatencySimulator;
import com.example.review.dto.response.LatestReviewResponse;
import com.example.review.dto.response.ProductReviewSummaryResponse;
import com.example.review.entity.ProductReview;
import com.example.review.repository.ProductReviewRepository;
import com.example.review.service.ProductReviewService;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;

import java.util.List;

@Component
@RequiredArgsConstructor
public class GetProductReviewSummaryCommand
    implements Command<GetProductReviewSummaryCommandRequest, ProductReviewSummaryResponse> {

    private final DatabaseLatencySimulator databaseLatencySimulator;

    private final ProductReviewRepository repository;

    private final ProductReviewService service;

    @Override
    public ProductReviewSummaryResponse execute(GetProductReviewSummaryCommandRequest request) {
        List<ProductReview> reviews = getByProductId(request);
        List<LatestReviewResponse> latestReviews = getLatestReviews(request);
        return new ProductReviewSummaryResponse(
            request.productId(),
            reviews.size(),
            service.calculateAverageRating(reviews),
            service.calculateRatingDistribution(reviews),
            latestReviews
        );
    }

    private List<ProductReview> getByProductId(GetProductReviewSummaryCommandRequest request) {
        databaseLatencySimulator.delay();
        return repository.findByProductId(request.productId());
    }

    private List<LatestReviewResponse> getLatestReviews(GetProductReviewSummaryCommandRequest request) {
        databaseLatencySimulator.delay();
        return repository.findTop5ByProductIdOrderByCreatedAtDesc(request.productId())
            .stream()
            .map(this::toLatestReviewResponse)
            .toList();
    }

    private LatestReviewResponse toLatestReviewResponse(ProductReview review) {
        return new LatestReviewResponse(
            review.getId(),
            review.getCustomerId(),
            review.getRating(),
            review.getComment(),
            review.getCreatedAt()
        );
    }

}
