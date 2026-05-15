package com.example.review.dto.response;

import java.util.List;

public record SellerProductsResponse(
    String sellerId,
    List<String> productIds
) {

}
