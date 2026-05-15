package http_router

import (
	"example/product-review-api-command-go/internal/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func InitializeRouter(httpHandler *handler.Handler) *mux.Router {
	router := mux.NewRouter()
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	apiV1.HandleFunc("/products/{productId}/reviews", httpHandler.ProductReviewHandler.CreateReview).Methods(http.MethodPost)
	apiV1.HandleFunc("/products/{productId}/review-summary", httpHandler.ProductReviewHandler.GetProductReviewSummary).Methods(http.MethodGet)
	apiV1.HandleFunc("/sellers/{sellerId}/review-analytics", httpHandler.ProductReviewHandler.GetSellerReviewAnalytics).Methods(http.MethodGet)
	apiV1.HandleFunc("/products/{productId}/reviews", httpHandler.ProductReviewHandler.GetProductReviews).Methods(http.MethodGet)
	apiV1.HandleFunc("/customers/{customerId}/reviews", httpHandler.ProductReviewHandler.GetCustomerReviews).Methods(http.MethodGet)
	apiV1.HandleFunc("/reviews/recent", httpHandler.ProductReviewHandler.GetRecentReviews).Methods(http.MethodGet)
	apiV1.HandleFunc("/products/{productId}/review-analytics/daily", httpHandler.ProductReviewHandler.GetDailyProductReviewAnalytics).Methods(http.MethodGet)

	return router
}
