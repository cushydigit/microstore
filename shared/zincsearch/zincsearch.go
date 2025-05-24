package zincsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/cushydigit/microstore/shared/types"
)

func (c *Client) healthCheck(ctx context.Context) error {
	req, err := c.newRequest(ctx, http.MethodGet, "/", nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) IndexProduct(ctx context.Context, index string, product *types.Product) error {
	url := fmt.Sprintf("/api/%s/_doc/%d", index, product.ID)
	req, err := c.newRequest(ctx, http.MethodPut, url, product)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to index product: status %d, body %s", resp.StatusCode, string(body))
	}

	return nil
}

// indexed but not able to delete a record with id (id not found)
func (c *Client) IndexBulkProductV2(ctx context.Context, index string, products []*types.Product) error {
	records := make([]map[string]any, 0, len(products))
	for _, p := range products {
		log.Printf("the id of prodcut: %d\n", p.ID)
		record := map[string]any{
			"_id":         fmt.Sprintf("%d", p.ID),
			"name":        p.Name,
			"description": p.Description,
			"price":       p.Price,
			"stock":       p.Stock,
			"@timestamp":  time.Now().UTC(),
		}
		records = append(records, record)
	}
	payload := map[string]any{
		"index":   index,
		"records": products,
	}
	req, err := c.newRequest(ctx, http.MethodPost, "/api/_bulkv2", payload)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bulk v2 indexing failed: status %d, body: %s", resp.StatusCode, respBody)
	}

	return nil

}

// not indexed
func (c *Client) IndexBulkProductv1(ctx context.Context, index string, products []*types.Product) error {
	var buffer bytes.Buffer
	for _, p := range products {
		meta := map[string]map[string]string{
			"index": {
				"_index": index,
				"_id":    fmt.Sprintf("%d", p.ID),
			},
		}
		metaJSON, err := json.Marshal(meta)
		if err != nil {
			return err
		}
		buffer.Write(metaJSON)
		buffer.WriteByte('\n')

		dataJSON, err := json.Marshal(p)
		if err != nil {
			return err
		}
		buffer.Write(dataJSON)
		buffer.WriteByte('\n')

	}
	req, err := c.newRequest(ctx, http.MethodPost, "/api/_bulk", buffer)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-ndjson")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bulk v indexing failed: status %d, body: %s", resp.StatusCode, respBody)
	}

	return nil

}
func (c *Client) IndexBulkProduct(ctx context.Context, index string, products []*types.Product) error {
	for _, p := range products {
		if err := c.IndexProduct(ctx, index, p); err != nil {
			return fmt.Errorf("failed to index in many products: %v", err)
		}
	}
	return nil
}

func (c *Client) DeleteProduct(ctx context.Context, index string, id int64) error {
	url := fmt.Sprintf("/api/%s/_doc/%d", index, id)
	req, err := c.newRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete indexed product: status %d, body %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) DeleteAllProducts(ctx context.Context, index string) error {
	url := fmt.Sprintf("/api/index/%s", index)

	req, err := c.newRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete index and all documents: status %d, body: %s", resp.StatusCode, respBody)
	}
	// create products index
	if err := c.createIndex(ctx, index); err != nil {
		return fmt.Errorf("failed to create index after deleting the index: %v", err)
	}

	return nil
}

func (c *Client) SearchProduct(ctx context.Context, index string, query string) ([]*types.Product, error) {
	reqBody := map[string]any{
		"search_type": "match",
		"query": map[string]any{
			"term": query,
		},
		"fields":     []string{"name", "description"},
		"from":       0,
		"max_result": 10,
	}
	url := fmt.Sprintf("/api/%s/_search", index)
	req, err := c.newRequest(ctx, http.MethodPost, url, reqBody)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res struct {
		Hits struct {
			Hits []struct {
				Source *types.Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	var products []*types.Product
	for _, hit := range res.Hits.Hits {
		products = append(products, hit.Source)
	}

	return products, nil

}

func (c *Client) createIndex(ctx context.Context, index string) error {
	reqBody := map[string]any{
		"name":         index,
		"storage_type": "disk",
	}
	req, err := c.newRequest(ctx, http.MethodPost, "/api/index", reqBody)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create index: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) indexExists(ctx context.Context, index string) (bool, error) {
	url := fmt.Sprintf("/api/index/%s", index)
	req, err := c.newRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	body, _ := io.ReadAll(resp.Body)
	return false, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
}
