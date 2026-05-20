package com.example.review.command.productreview;

import com.example.review.command.Command;
import com.example.review.dto.command.request.GetListRecentReviewsCommandRequest;
import com.example.review.config.DatabaseLatencySimulator;
import com.example.review.dto.response.ReviewDetailResponse;
import com.example.review.entity.ProductReview;
import com.example.review.repository.ProductReviewRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.PageRequest;
import org.springframework.stereotype.Component;

import java.util.List;

@Component
@RequiredArgsConstructor
public class GetListRecentReviewsCommand implements Command<GetListRecentReviewsCommandRequest, List<ReviewDetailResponse>> {

    private static final int MAX_RECENT_REVIEW_LIMIT = 50;

    private static final int MIN_RECENT_REVIEW_LIMIT = 1;

    private final DatabaseLatencySimulator databaseLatencySimulator;

    private final ProductReviewRepository repository;

    @Override
    public List<ReviewDetailResponse> execute(GetListRecentReviewsCommandRequest request) {
        int normalizedLimit = getNormalizeRecentReviewLimit(request.limit());
        databaseLatencySimulator.delay();
        return repository.findAllByOrderByCreatedAtDesc(PageRequest.of(0, normalizedLimit))
            .stream()
            .map(this::toReviewDetailResponse)
            .toList();
    }

    public int getNormalizeRecentReviewLimit(int limit) {
        if (limit < MIN_RECENT_REVIEW_LIMIT) {
            return MIN_RECENT_REVIEW_LIMIT;
        }
        return Math.min(limit, MAX_RECENT_REVIEW_LIMIT);
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
