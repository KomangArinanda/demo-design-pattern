package client

import "example/product-review-api-command-go/internal/dto/response"

type ProductClient interface {
	GetProductIDsBySellerID(sellerID string) response.SellerProductsResponse
}

type productClient struct {
	sellerProducts map[string][]string
}

func NewProductClient() ProductClient {
	return &productClient{
		sellerProducts: map[string][]string{
			"SELLER-001": {"PROD-001", "PROD-002", "PROD-003"},
			"SELLER-002": {"PROD-004", "PROD-005"},
		},
	}
}

func (c *productClient) GetProductIDsBySellerID(sellerID string) response.SellerProductsResponse {
	return response.SellerProductsResponse{
		SellerID:   sellerID,
		ProductIDs: c.sellerProducts[sellerID],
	}
}
