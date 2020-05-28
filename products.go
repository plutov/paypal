package paypal

import (
	"fmt"
	"net/http"
)

type (
	ProductType string
	ProductCategory string

	// Product struct
	Product struct {
		ID          string              `json:"id,omitempty"`
		Name        string              `json:"name,omitempty"`
		Description string `json:"description"`
		ProductType ProductType `json:"description"`
		ImageUrl string `json:"image_url"`
		HomeUrl string `json:"home_url"`
	}

	CreateProductResponse struct {
		Product
		CreateTime string `json:"create_time"`
		UpdateTime string `json:"update_time"`
		Links []Link `json:"links"`
	}
)

const (
	PRODUCT_TYPE_PHYSICAL ProductType = "PHYSICAL"
	PRODUCT_TYPE_DIGITAL  ProductType = "DIGITAL"
	PRODUCT_TYPE_SERVICE  ProductType = "SERVICE"

	PRODUCT_CATEGORY_SOFTWARE ProductCategory = "software"
)

// CreateProduct creates a product
// Endpoint: POST /v1/catalogs/products
func (c *Client) CreateProduct(product Product) (*CreateProductResponse, error) {
	req, err := c.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", c.APIBase, "/v1/catalogs/products"), product)
	response := &CreateProductResponse{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}