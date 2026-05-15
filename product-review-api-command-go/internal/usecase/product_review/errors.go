package product_review

import "errors"

var (
	ErrDuplicateReview = errors.New("Duplicate review")
	ErrInvalidOrder    = errors.New("Customer did not purchase this product")
	ErrInvalidRating   = errors.New("Rating must be between 1 and 5")
	ErrInvalidMonth    = errors.New("Month must be between 1 and 12")
)
