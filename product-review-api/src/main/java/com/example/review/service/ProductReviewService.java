package com.example.review.service;

import com.example.review.client.OrderClient;
import com.example.review.client.ProductClient;
import com.example.review.config.DatabaseLatencySimulator;
import com.example.review.dto.request.CreateReviewRequest;
import com.example.review.dto.response.*;
import com.example.review.entity.ProductReview;
import com.example.review.exception.DuplicateReviewException;
import com.example.review.exception.InvalidOrderException;
import com.example.review.repository.ProductReviewRepository;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.data.domain.PageRequest;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;
import java.math.RoundingMode;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.YearMonth;
import java.util.*;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
public class ProductReviewService {

    private static final int MAX_RECENT_REVIEW_LIMIT = 50;

    private static final int MIN_RECENT_REVIEW_LIMIT = 1;

    private static final Logger log = LoggerFactory.getLogger(ProductReviewService.class);

    private final DatabaseLatencySimulator databaseLatencySimulator;

    private final OrderClient orderClient;

    private final ProductClient productClient;

    private final ProductReviewRepository repository;

    private double calculateAverageRating(List<ProductReview> reviews) {
        if (reviews.isEmpty()) {
            return 0.0;
        }
        double average = reviews.stream()
            .mapToInt(ProductReview::getRating)
            .average()
            .orElse(0.0);
        return round(average);
    }

    private double calculatePercentage(long numerator, long denominator) {
        if (denominator == 0) {
            return 0.0;
        }
        return round((double) numerator * 100 / denominator);
    }

    private Map<Integer, Long> calculateRatingDistribution(List<ProductReview> reviews) {
        Map<Integer, Long> distribution = new LinkedHashMap<>();
        for (int rating = 1; rating <= 5; rating++) {
            distribution.put(rating, 0L);
        }
        reviews.forEach(review -> distribution.computeIfPresent(review.getRating(), (rating, count) -> count + 1));
        return distribution;
    }

    @Transactional
    public CreateReviewResponse createReview(String productId, CreateReviewRequest request) {
        log.info(
            "Creating review for productId={}, customerId={}, orderId={}",
            productId,
            request.customerId(),
            request.orderId()
        );

        OrderValidationResponse orderValidation = orderClient.validateOrder(
            request.customerId(),
            request.orderId(),
            productId
        );

        if (!orderValidation.valid() || !"COMPLETED".equals(orderValidation.orderStatus())) {
            throw new InvalidOrderException();
        }

        databaseLatencySimulator.delay();
        if (repository.existsByProductIdAndCustomerIdAndOrderId(
            productId,
            request.customerId(),
            request.orderId()
        )) {
            throw new DuplicateReviewException();
        }

        databaseLatencySimulator.delay();
        ProductReview savedReview = repository.save(new ProductReview(
            productId,
            request.customerId(),
            request.orderId(),
            request.rating(),
            request.comment()
        ));

        return toCreateReviewResponse(savedReview);
    }

    @Transactional(readOnly = true)
    public List<ReviewDetailResponse> getCustomerReviews(String customerId) {
        log.info("Fetching reviews for customerId={}", customerId);
        databaseLatencySimulator.delay();
        return repository.findByCustomerIdOrderByCreatedAtDesc(customerId)
            .stream()
            .map(this::toReviewDetailResponse)
            .toList();
    }

    @Transactional(readOnly = true)
    public DailyProductReviewAnalyticsResponse getDailyProductReviewAnalytics(
        String productId,
        int month,
        int year
    ) {
        YearMonth requestedMonth = YearMonth.of(year, month);
        LocalDateTime startDateTime = requestedMonth.atDay(1).atStartOfDay();
        LocalDateTime endDateTime = requestedMonth.plusMonths(1).atDay(1).atStartOfDay();

        databaseLatencySimulator.delay();
        List<ProductReview> monthlyReviews =
            repository.findByProductIdAndCreatedAtGreaterThanEqualAndCreatedAtLessThan(
                productId,
                startDateTime,
                endDateTime
            );

        // Kept separate intentionally so the service remains the aggregation owner
        // while still showing a second repository interaction for the same use case.
        databaseLatencySimulator.delay();
        List<ProductReview> monthlyReviewsByRecency =
            repository.findByProductIdAndCreatedAtGreaterThanEqualAndCreatedAtLessThanOrderByCreatedAtDesc(
                productId,
                startDateTime,
                endDateTime
            );

        Map<LocalDate, List<ProductReview>> reviewsByDate = monthlyReviews.stream()
            .collect(Collectors.groupingBy(
                review -> review.getCreatedAt().toLocalDate(),
                TreeMap::new,
                Collectors.toList()
            ));

        List<DailyReviewSummaryResponse> dailySummaries = new ArrayList<>();
        for (Map.Entry<LocalDate, List<ProductReview>> entry : reviewsByDate.entrySet()) {
            List<ProductReview> dailyReviews = entry.getValue();
            dailySummaries.add(new DailyReviewSummaryResponse(
                entry.getKey(),
                dailyReviews.size(),
                calculateAverageRating(dailyReviews),
                calculateRatingDistribution(dailyReviews)
            ));
        }

        log.info(
            "Built daily analytics for productId={}, month={}, year={}, totalReviews={}, latestMonthlyReviewCount={}",
            productId,
            month,
            year,
            monthlyReviews.size(),
            monthlyReviewsByRecency.size()
        );

        return new DailyProductReviewAnalyticsResponse(
            productId,
            month,
            year,
            monthlyReviews.size(),
            calculateAverageRating(monthlyReviews),
            calculateRatingDistribution(monthlyReviews),
            dailySummaries
        );
    }

    @Transactional(readOnly = true)
    public ProductReviewSummaryResponse getProductReviewSummary(String productId) {
        databaseLatencySimulator.delay();
        List<ProductReview> reviews = repository.findByProductId(productId);
        databaseLatencySimulator.delay();
        return new ProductReviewSummaryResponse(
            productId,
            reviews.size(),
            calculateAverageRating(reviews),
            calculateRatingDistribution(reviews),
            repository.findTop5ByProductIdOrderByCreatedAtDesc(productId)
                .stream()
                .map(this::toLatestReviewResponse)
                .toList()
        );
    }

    @Transactional(readOnly = true)
    public List<ReviewDetailResponse> getProductReviews(String productId) {
        log.info("Fetching reviews for productId={}", productId);
        databaseLatencySimulator.delay();
        return repository.findByProductIdOrderByCreatedAtDesc(productId)
            .stream()
            .map(this::toReviewDetailResponse)
            .toList();
    }

    @Transactional(readOnly = true)
    public List<ReviewDetailResponse> getRecentReviews(int limit) {
        int normalizedLimit = normalizeRecentReviewLimit(limit);
        log.info("Fetching recent reviews with limit={}", normalizedLimit);
        databaseLatencySimulator.delay();
        return repository.findAllByOrderByCreatedAtDesc(PageRequest.of(0, normalizedLimit))
            .stream()
            .map(this::toReviewDetailResponse)
            .toList();
    }

    @Transactional(readOnly = true)
    public SellerReviewAnalyticsResponse getSellerReviewAnalytics(String sellerId) {
        SellerProductsResponse sellerProducts = productClient.getProductIdsBySellerId(sellerId);
        databaseLatencySimulator.delay();
        List<ProductReview> reviews = repository.findByProductIdIn(sellerProducts.productIds());
        long negativeReviewCount = reviews.stream()
            .filter(review -> review.getRating() <= 2)
            .count();

        Map<String, Long> negativeCountByProduct = reviews.stream()
            .filter(review -> review.getRating() <= 2)
            .collect(Collectors.groupingBy(ProductReview::getProductId, Collectors.counting()));

        List<TopComplainedProductResponse> topComplainedProducts = negativeCountByProduct.entrySet()
            .stream()
            .sorted(Map.Entry.<String, Long>comparingByValue(Comparator.reverseOrder())
                .thenComparing(Map.Entry.comparingByKey()))
            .limit(3)
            .map(entry -> new TopComplainedProductResponse(entry.getKey(), entry.getValue()))
            .toList();

        return new SellerReviewAnalyticsResponse(
            sellerId,
            reviews.size(),
            calculateAverageRating(reviews),
            negativeReviewCount,
            calculatePercentage(negativeReviewCount, reviews.size()),
            calculateRatingDistribution(reviews),
            topComplainedProducts
        );
    }

    private int normalizeRecentReviewLimit(int limit) {
        if (limit < MIN_RECENT_REVIEW_LIMIT) {
            return MIN_RECENT_REVIEW_LIMIT;
        }
        return Math.min(limit, MAX_RECENT_REVIEW_LIMIT);
    }

    private double round(double value) {
        return BigDecimal.valueOf(value)
            .setScale(1, RoundingMode.HALF_UP)
            .doubleValue();
    }

    private CreateReviewResponse toCreateReviewResponse(ProductReview review) {
        return new CreateReviewResponse(
            review.getId(),
            review.getProductId(),
            review.getCustomerId(),
            review.getOrderId(),
            review.getRating(),
            review.getComment(),
            review.getCreatedAt()
        );
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
