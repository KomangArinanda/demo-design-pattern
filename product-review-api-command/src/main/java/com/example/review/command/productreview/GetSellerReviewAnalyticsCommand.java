package com.example.review.command.productreview;

import com.example.review.client.ProductClient;
import com.example.review.command.Command;
import com.example.review.dto.command.request.GetSellerReviewAnalyticsCommandRequest;
import com.example.review.config.DatabaseLatencySimulator;
import com.example.review.dto.response.SellerProductsResponse;
import com.example.review.dto.response.SellerReviewAnalyticsResponse;
import com.example.review.dto.response.TopComplainedProductResponse;
import com.example.review.entity.ProductReview;
import com.example.review.repository.ProductReviewRepository;
import com.example.review.service.ProductReviewService;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;

import java.util.Comparator;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

import static java.lang.Math.round;

@Component
@RequiredArgsConstructor
public class GetSellerReviewAnalyticsCommand
    implements Command<GetSellerReviewAnalyticsCommandRequest, SellerReviewAnalyticsResponse> {

    private final DatabaseLatencySimulator databaseLatencySimulator;

    private final ProductClient productClient;

    private final ProductReviewRepository repository;

    private final ProductReviewService service;

    private double calculatePercentage(long numerator, long denominator) {
        if (denominator == 0) {
            return 0.0;
        }
        return round((double) numerator * 100 / denominator);
    }

    @Override
    public SellerReviewAnalyticsResponse execute(GetSellerReviewAnalyticsCommandRequest request) {
        SellerProductsResponse sellerProducts = productClient.getProductIdsBySellerId(request.sellerId());
        List<ProductReview> reviews = getByProductIdIn(sellerProducts);
        long negativeReviewCount = getNegativeReviewCount(reviews);
        Map<String, Long> negativeCountByProduct = getNegativeCountByProduct(reviews);
        List<TopComplainedProductResponse> topComplainedProducts = getTopComplainedProducts(negativeCountByProduct);

        return new SellerReviewAnalyticsResponse(
            request.sellerId(),
            reviews.size(),
            service.calculateAverageRating(reviews),
            negativeReviewCount,
            calculatePercentage(negativeReviewCount, reviews.size()),
            service.calculateRatingDistribution(reviews),
            topComplainedProducts
        );
    }

    private List<ProductReview> getByProductIdIn(SellerProductsResponse sellerProducts) {
        databaseLatencySimulator.delay();
        return repository.findByProductIdIn(sellerProducts.productIds());
    }

    private Map<String, Long> getNegativeCountByProduct(List<ProductReview> reviews) {
        return reviews.stream()
            .filter(review -> review.getRating() <= 2)
            .collect(Collectors.groupingBy(ProductReview::getProductId, Collectors.counting()));
    }

    private long getNegativeReviewCount(List<ProductReview> reviews) {
        return reviews.stream()
            .filter(review -> review.getRating() <= 2)
            .count();
    }

    private List<TopComplainedProductResponse> getTopComplainedProducts(Map<String, Long> negativeCountByProduct) {
        return negativeCountByProduct.entrySet()
            .stream()
            .sorted(Map.Entry.<String, Long>comparingByValue(Comparator.reverseOrder())
                .thenComparing(Map.Entry.comparingByKey()))
            .limit(3)
            .map(entry -> new TopComplainedProductResponse(entry.getKey(), entry.getValue()))
            .toList();
    }

}
