package usecase

import (
	"context"

	"example/product-review-api-command-go/internal/shared/appctx"
)

type Usecase interface {
	Execute(ctx context.Context, input any) appctx.Response
}
