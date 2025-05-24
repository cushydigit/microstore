package zincsearch

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/cushydigit/microstore/shared/types"
)

type Client struct {
	BaseUrl    string
	Username   string
	Password   string
	HTTPClient *http.Client
}

var (
	instance *Client
	once     sync.Once
)

func Init(baseUrl, username, password, index string) *Client {
	once.Do(func() {
		instance = &Client{
			BaseUrl:  baseUrl,
			Username: username,
			Password: password,
			HTTPClient: &http.Client{
				Timeout: 10 * time.Second,
			},
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// check the connection
	if err := instance.healthCheck(ctx); err != nil {
		log.Fatalf("failed to connect to zincsearch: %v", err)
		return nil
	}

	// check if index is exists
	ok, err := instance.indexExists(ctx, index)
	if err != nil {
		log.Fatalf("failed to check if index is exists: %v", err)
		return nil
	}

	if !ok {
		if err := instance.createIndex(ctx, index); err != nil {
			log.Fatalf("failed to create not existed index : %v", err)
			return nil
		}
	}

	return instance
}

func GetInstance() *Client {
	return instance
}

func (c *Client) healthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseUrl+"/", nil)
	if err != nil {
		return err
	}
	resp, err := c.HTTPClient.Do(req)
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
	url := fmt.Sprintf("%s/api/%s/_doc/%d", c.BaseUrl, index, product.ID)
	log.Printf("url: %s", url)

	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
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

func (c *Client) IndexBulkProduct(ctx context.Context, index string, products []*types.Product) error {
	payload := map[string]any{
		"index":   index,
		"records": products,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal bulk payload: %w", err)
	}

	url := fmt.Sprintf("%s/api/_bulkv2", c.BaseUrl)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Basic "+basicAuth(c.Username, c.Password))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
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

func (c *Client) DeleteProduct(ctx context.Context, index string, id int64) error {
	url := fmt.Sprintf("%s/api/%s/_doc/%d", c.BaseUrl, index, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
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

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (c *Client) SearchProduct(ctx context.Context, query string) ([]*types.Product, error) {
	reqBody := map[string]any{
		"search_type": "match",
		"query": map[string]any{
			"term": query,
		},
		"fields":     []string{"name", "description"},
		"from":       0,
		"max_result": 10,
	}
	body, _ := json.Marshal(reqBody)
	url := c.BaseUrl + "/api/products/_search"
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	req.Header.Set("Authorization", "Basic "+basicAuth(c.Username, c.Password))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// FOR DEBUGGING
	// rawBody, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to read response body: %v", err)
	// }

	// pretty-print the JSON
	// var pretty bytes.Buffer
	// if err := json.Indent(&pretty, rawBody, "", "  "); err != nil {
	// 	fmt.Println("Raw body (not valid JSON):")
	// 	fmt.Println(string(rawBody))
	// } else {
	// 	fmt.Println("Pretty JSON response: ")
	// 	fmt.Println(pretty.String())
	// }

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

	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseUrl+"/api/index", bytes.NewReader(body))
	if err != nil {
		return nil
	}

	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
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
	url := fmt.Sprintf("%s/api/index/%s", c.BaseUrl, index)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	req.SetBasicAuth(c.Username, c.Password)
	resp, err := c.HTTPClient.Do(req)
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
