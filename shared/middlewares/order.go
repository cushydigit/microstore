package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/cushydigit/microstore/shared/helpers"
	"github.com/cushydigit/microstore/shared/types"
)

func ValidateCreateOrder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateOrderRequest
		if err := helpers.ReadJSON(w, r, &req); err != nil {
			helpers.ErrorJSON(w, errors.New("invalid request body"))
			return
		}

		// basic validation
		if len(req.Items) == 0 {
			helpers.ErrorJSON(w, errors.New("order must contain at least one item"))
			return
		}
		for _, item := range req.Items {
			if item.ProductID <= 0 || item.Quantity <= 0 {
				helpers.ErrorJSON(w, errors.New("invalid item values"))
				return
			}
		}

		ctx := context.WithValue(r.Context(), types.CreateOrderRequestKey, req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
