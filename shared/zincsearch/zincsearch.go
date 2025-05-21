package zincsearch

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
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

func Init(baseUrl, username, password string) *Client {
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

	if err := instance.HealthCheck(ctx); err != nil {
		log.Fatalf("failed to connect to zincsearch: %v", err)
		return nil
	}

	return instance
}

func GetInstance() *Client {
	return instance
}

func (c *Client) HealthCheck(ctx context.Context) error {
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
