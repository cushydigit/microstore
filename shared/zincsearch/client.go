package zincsearch

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Client struct {
	baseUrl    string
	authHeader string
	httpClient *http.Client
}

var (
	instance *Client
	once     sync.Once
)

func GetInstanct() *Client {
	return instance
}

func Init(baseUrl, username, password, index string) *Client {
	once.Do(func() {
		auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
		instance = &Client{
			baseUrl:    strings.TrimSuffix(baseUrl, "/"),
			authHeader: "Basic " + auth,
			httpClient: &http.Client{
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

func (c *Client) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseUrl+path, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.authHeader)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
