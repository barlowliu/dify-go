package dify-go

import (
    "fmt"
    "time"

    "github.com/hashicorp/go-retryablehttp"
)

// Client is the main struct for interacting with the Dify API.
type Client struct {
    BaseURL    string
    APIKey     string
    HTTPClient *retryablehttp.Client
}

// NewClient initializes and returns a new Dify API client.
func NewClient(baseURL, apiKey string) *Client {
    client := retryablehttp.NewClient()
    client.RetryMax = 3
    client.RetryWaitMin = 500 * time.Millisecond
    client.RetryWaitMax = 2 * time.Second

    return &Client{
        BaseURL:    baseURL,
        APIKey:     apiKey,
        HTTPClient: client,
    }
}

// SetTimeout allows setting a custom timeout for the HTTP client.
func (c *Client) SetTimeout(timeout time.Duration) {
    c.HTTPClient.HTTPClient.Timeout = timeout
}

// buildURL constructs the full API endpoint URL.
func (c *Client) buildURL(endpoint string) string {
    return fmt.Sprintf("%s%s", c.BaseURL, endpoint)
}

// addHeaders adds the necessary headers to the request.
func (c *Client) addHeaders(req *retryablehttp.Request) {
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
    req.Header.Set("Content-Type", "application/json")
}

