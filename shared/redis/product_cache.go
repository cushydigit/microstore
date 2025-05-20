package redis

import (
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

func GetProductFromCache(id int64) (*types.Product, bool, error) {
	data, err := Client.Get(Ctx, productKey(id)).Result()
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

func SetProductToCache(product *types.Product) error {
	jsonData, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return Client.Set(Ctx, productKey(product.ID), jsonData, productTTL).Err()
}

func DeleteProductFromCache(id int64) error {
	return Client.Del(Ctx, productKey(id)).Err()
}
