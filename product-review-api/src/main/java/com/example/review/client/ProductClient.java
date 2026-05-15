package com.example.review.client;

import com.example.review.dto.response.SellerProductsResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;

import java.util.List;
import java.util.Map;

@Component
public class ProductClient {

    private static final Map<String, List<String>> SELLER_PRODUCTS = Map.of(
        "SELLER-001", List.of("PROD-001", "PROD-002", "PROD-003"),
        "SELLER-002", List.of("PROD-004", "PROD-005")
    );

    private static final Logger log = LoggerFactory.getLogger(ProductClient.class);

    public SellerProductsResponse getProductIdsBySellerId(String sellerId) {
        log.info("Fetching product list for sellerId={}", sellerId);
        return new SellerProductsResponse(sellerId, SELLER_PRODUCTS.getOrDefault(sellerId, List.of()));
    }
}
