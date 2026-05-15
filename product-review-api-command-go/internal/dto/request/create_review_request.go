package request

type CreateReviewRequest struct {
	CustomerID string `json:"customerId"`
	OrderID    string `json:"orderId"`
	Rating     int    `json:"rating"`
	Comment    string `json:"comment"`
}
