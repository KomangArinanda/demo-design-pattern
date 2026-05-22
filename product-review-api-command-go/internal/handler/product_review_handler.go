package handler

import (
	"encoding/json"
	"example/product-review-api-command-go/internal/dto/request"
	"example/product-review-api-command-go/internal/shared/appctx"
	product_review "example/product-review-api-command-go/internal/usecase/product_review"
	"example/product-review-api-command-go/internal/usecase/product_review/get_daily_analytics"
	"example/product-review-api-command-go/internal/usecase/product_review/get_list"
	"example/product-review-api-command-go/internal/usecase/product_review/get_list_customer_reviews"
	"example/product-review-api-command-go/internal/usecase/product_review/get_list_recent_reviews"
	"example/product-review-api-command-go/internal/usecase/product_review/get_seller_review_analytics"
	"example/product-review-api-command-go/internal/usecase/product_review/get_summary"
	"example/product-review-api-command-go/internal/usecase/product_review/save"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ProductReviewHandler interface {
	CreateReview(w http.ResponseWriter, r *http.Request)
	GetProductReviewSummary(w http.ResponseWriter, r *http.Request)
	GetSellerReviewAnalytics(w http.ResponseWriter, r *http.Request)
	GetProductReviews(w http.ResponseWriter, r *http.Request)
	GetCustomerReviews(w http.ResponseWriter, r *http.Request)
	GetRecentReviews(w http.ResponseWriter, r *http.Request)
	GetDailyProductReviewAnalytics(w http.ResponseWriter, r *http.Request)
}

type productReviewHandler struct {
	usecases *product_review.ProductReview
}

func NewProductReviewHandler(usecases *product_review.ProductReview) ProductReviewHandler {
	return &productReviewHandler{usecases: usecases}
}

func (h *productReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	var body request.CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	requestModel := save.Request{
		ProductID: mux.Vars(r)["productId"],
		Request:   body,
	}
	writeAppResponse(w, h.usecases.Save.Execute(r.Context(), requestModel))
}

func (h *productReviewHandler) GetProductReviewSummary(w http.ResponseWriter, r *http.Request) {
	writeAppResponse(w, h.usecases.GetSummary.Execute(r.Context(), get_summary.GetSummaryRequest{
		ProductID: mux.Vars(r)["productId"],
	}))
}

func (h *productReviewHandler) GetSellerReviewAnalytics(w http.ResponseWriter, r *http.Request) {
	writeAppResponse(w, h.usecases.GetSellerReviewAnalytics.Execute(r.Context(), get_seller_review_analytics.GetSellerReviewAnalyticsRequest{
		SellerID: mux.Vars(r)["sellerId"],
	}))
}

func (h *productReviewHandler) GetProductReviews(w http.ResponseWriter, r *http.Request) {
	writeAppResponse(w, h.usecases.GetList.Execute(r.Context(), get_list.GetListRequest{
		ProductID: mux.Vars(r)["productId"],
	}))
}

func (h *productReviewHandler) GetCustomerReviews(w http.ResponseWriter, r *http.Request) {
	writeAppResponse(w, h.usecases.GetListCustomerReviews.Execute(r.Context(), get_list_customer_reviews.GetListCustomerReviewsRequest{
		CustomerID: mux.Vars(r)["customerId"],
	}))
}

func (h *productReviewHandler) GetRecentReviews(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 10
	}
	writeAppResponse(w, h.usecases.GetListRecentReviews.Execute(r.Context(), get_list_recent_reviews.GetListRecentReviewsRequest{
		Limit: limit,
	}))
}

func (h *productReviewHandler) GetDailyProductReviewAnalytics(w http.ResponseWriter, r *http.Request) {
	month, err := strconv.Atoi(r.URL.Query().Get("month"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Month must be between 1 and 12")
		return
	}
	year, err := strconv.Atoi(r.URL.Query().Get("year"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Year is required")
		return
	}
	requestModel := get_daily_analytics.GetDailyAnalyticsRequest{
		ProductID: mux.Vars(r)["productId"],
		Month:     month,
		Year:      year,
	}
	writeAppResponse(w, h.usecases.GetDailyProductReviewAnalytics.Execute(r.Context(), requestModel))
}

func writeAppResponse(w http.ResponseWriter, res appctx.Response) {
	if res.Code >= http.StatusBadRequest {
		writeError(w, res.Code, res.Message)
		return
	}
	writeJSON(w, res.Code, res.Data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"message": message})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
