package com.example.review.command.productreview;

import com.example.review.command.Command;
import com.example.review.dto.command.request.GetDailyProductReviewAnalyticsCommandRequest;
import com.example.review.config.DatabaseLatencySimulator;
import com.example.review.dto.response.DailyProductReviewAnalyticsResponse;
import com.example.review.dto.response.DailyReviewSummaryResponse;
import com.example.review.entity.ProductReview;
import com.example.review.repository.ProductReviewRepository;
import com.example.review.service.ProductReviewService;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.YearMonth;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.TreeMap;
import java.util.stream.Collectors;

@Component
@RequiredArgsConstructor
public class GetDailyProductReviewAnalyticsCommand
    implements Command<GetDailyProductReviewAnalyticsCommandRequest, DailyProductReviewAnalyticsResponse> {

    private final DatabaseLatencySimulator databaseLatencySimulator;

    private final ProductReviewRepository repository;

    private final ProductReviewService service;

    @Override
    public DailyProductReviewAnalyticsResponse execute(GetDailyProductReviewAnalyticsCommandRequest request) {
        YearMonth requestedMonth = YearMonth.of(request.year(), request.month());
        LocalDateTime startDateTime = requestedMonth.atDay(1).atStartOfDay();
        LocalDateTime endDateTime = requestedMonth.plusMonths(1).atDay(1).atStartOfDay();

        databaseLatencySimulator.delay();
        List<ProductReview> monthlyReviews =
            repository.findByProductIdAndCreatedAtGreaterThanEqualAndCreatedAtLessThan(
                request.productId(),
                startDateTime,
                endDateTime
            );
        databaseLatencySimulator.delay();
        repository.findByProductIdAndCreatedAtGreaterThanEqualAndCreatedAtLessThanOrderByCreatedAtDesc(
            request.productId(),
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
                service.calculateAverageRating(dailyReviews),
                service.calculateRatingDistribution(dailyReviews)
            ));
        }

        return new DailyProductReviewAnalyticsResponse(
            request.productId(),
            request.month(),
            request.year(),
            monthlyReviews.size(),
            service.calculateAverageRating(monthlyReviews),
            service.calculateRatingDistribution(monthlyReviews),
            dailySummaries
        );
    }
}
