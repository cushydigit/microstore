package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/cushydigit/microstore/shared/helpers"
	"github.com/cushydigit/microstore/shared/types"
)

func ValidateCreateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var product types.Product
		if err := helpers.ReadJSON(w, r, &product); err != nil {
			helpers.ErrorJSON(w, errors.New("invalid request body"))
			return
		}

		// basic validation
		if product.Name == "" || product.Price <= 0 || product.Stock < 0 {
			helpers.ErrorJSON(w, errors.New("invalid product fields"))
			return
		}

		// inject validated product into context
		ctx := context.WithValue(r.Context(), types.ProductKey, product)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
