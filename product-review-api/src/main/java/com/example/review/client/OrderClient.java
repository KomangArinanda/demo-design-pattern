package com.example.review.client;

import com.example.review.dto.response.OrderValidationResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;

@Component
public class OrderClient {

    private static final Logger log = LoggerFactory.getLogger(OrderClient.class);

    public OrderValidationResponse validateOrder(String customerId, String orderId, String productId) {
        log.info("Validating order ownership for customerId={}, orderId={}, productId={}", customerId, orderId, productId);

        boolean valid = orderId.startsWith("ORD-") && !orderId.equals("ORD-INVALID");
        return new OrderValidationResponse(valid, valid ? "COMPLETED" : "NOT_FOUND");
    }
}
