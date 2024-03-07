package dummyjson

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/The-Gleb/product_catalog/internal/domain/entity"
)

type productClient struct {
	url    string
	client http.Client
}

func NewProductClient(url string) *productClient {
	return &productClient{
		url:    url,
		client: *http.DefaultClient,
	}
}

func (c *productClient) GetNewProducts(ctx context.Context, offset int) ([]entity.AddOrUpdateProductDTO, error) {

	path := "/products?limit=10&skip=" + strconv.FormatInt(int64(offset), 10) + "&select=title,category"
	req, err := http.NewRequest("GET", c.url+path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	slog.Debug(string(body))

	var responseStruct struct {
		Products []entity.AddOrUpdateProductDTO
	}

	err = json.Unmarshal(body, &responseStruct)
	if err != nil {
		return nil, err
	}

	slog.Debug(fmt.Sprint(responseStruct.Products))

	return responseStruct.Products, nil

}
