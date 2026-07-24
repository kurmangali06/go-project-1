package product

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var ErrProductNotFound = errors.New("product not found")

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

type Product struct {
	Name  string
	Price uint32
}

func New(baseURL, token string) *Client {
	return &Client{
		baseURL:    baseURL,
		token:      token,
		httpClient: &http.Client{Timeout: 2 * time.Second},
	}
}

func (c *Client) GetProduct(ctx context.Context, sku int64) (Product, error) {
	reqBody := struct {
		Token string `json:"token"`
		Sku   uint32 `json:"sku"`
	}{
		Token: c.token,
		Sku:   uint32(sku),
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return Product{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/get_product", bytes.NewReader(body))
	if err != nil {
		return Product{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Product{}, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// ок, разбираем тело ниже
	case http.StatusNotFound:
		return Product{}, ErrProductNotFound
	default:
		return Product{}, fmt.Errorf("product service: status %d", resp.StatusCode)
	}

	var respBody struct {
		Name  string `json:"name"`
		Price uint32 `json:"price"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return Product{}, err
	}

	return Product{Name: respBody.Name, Price: respBody.Price}, nil
}
