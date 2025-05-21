package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cushydigit/microstore/shared/types"
	rds "github.com/redis/go-redis/v9"
)

const (
	productTTL = 5 * time.Minute
)

func productKey(id int64) string {
	return fmt.Sprintf("product:%s", strconv.FormatInt(id, 10))
}

func GetProductFromCache(ctx context.Context, id int64) (*types.Product, bool, error) {
	data, err := Client.Get(ctx, productKey(id)).Result()
	if err == rds.Nil {
		return nil, false, nil // Cache miss, not an error
	} else if err != nil {
		return nil, false, err // Redis connection error, etc
	}

	var product types.Product
	if err := json.Unmarshal([]byte(data), &product); err != nil {
		return nil, false, err
	}

	return &product, true, nil
}

func SetProductToCache(ctx context.Context, product *types.Product) error {
	jsonData, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return Client.Set(ctx, productKey(product.ID), jsonData, productTTL).Err()
}

func DeleteProductFromCache(ctx context.Context, id int64) error {
	return Client.Del(ctx, productKey(id)).Err()
}
