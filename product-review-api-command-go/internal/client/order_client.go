package client

import "example/product-review-api-command-go/internal/dto/response"

type OrderClient interface {
	ValidateOrder(customerID string, orderID string, productID string) response.OrderValidationResponse
}

type orderClient struct {
}

func NewOrderClient() OrderClient {
	return &orderClient{}
}

func (c *orderClient) ValidateOrder(customerID string, orderID string, productID string) response.OrderValidationResponse {
	valid := len(orderID) >= 4 && orderID[:4] == "ORD-" && orderID != "ORD-INVALID"
	status := "NOT_FOUND"
	if valid {
		status = "COMPLETED"
	}
	return response.OrderValidationResponse{
		Valid:       valid,
		OrderStatus: status,
	}
}
