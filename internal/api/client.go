package api

import (
	"net/http"
	"time"
)

// Client handles HTTP communication with the MYN backend.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetToken sets the Bearer token for authenticated requests.
func (c *Client) SetToken(token string) {
	c.Token = token
}
