package paypal

import (
	"fmt"
	"strconv"
)

// CreateProduct - Use this call to create a catalog product
// Endpoint: POST /v1/catalogs/products
// Type represents the product type. Indicates whether the product is physical or tangible goods, or a service. The allowed values are:
// -----------------------------
// | PHYSICAL | Physical goods |
// | DIGITAL  | Digital goods  |
// | SERVICE  | Service goods  |
// -----------------------------
func (c *Client) CreateProduct(product *CreateProduct) (*Product, error) {
	resp := &Product{}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/catalogs/products"), product)
	if err != nil {
		return nil, err
	}

	if err = c.SendWithBasicAuth(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// ListAllProducts lists all products
// Endpoint: GET /v1/catalogs/products
func (c *Client) ListAllProducts(params *ListProductsParams) (*ListProductsResponse, error) {
	resp := &ListProductsResponse{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/catalogs/products"), nil)
	q := req.URL.Query()
	q.Add("page_size", strconv.FormatUint(params.PageSize, 10))
	q.Add("page", strconv.FormatUint(params.Page, 10))
	q.Add("total_required", strconv.FormatBool(params.TotalRequired))
	if err != nil {
		return nil, err
	}

	if err = c.SendWithBasicAuth(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// ShowProduct shows details for a product by ID
// Endpoint: GET /v1/catalogs/products/{product_id}
func (c *Client) ShowProduct(productID string) (*Product, error) {
	resp := &Product{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/catalogs/products/"+productID), nil)
	if err != nil {
		return nil, err
	}

	if err = c.SendWithBasicAuth(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// UpdateProduct updates product by ID.
// You can update the following attributes and objects
//-----------------------------------------------
// | Attribute or Object | Operation			|
// ----------------------------------------------
// | /description         | add, replace, remove |
// | /category       	 | add, replace, remove |
// | /image_url			 | add, replace, remove |
// | /home_url			 | add, replace, remove |
// ----------------------------------------------
// Endpoint: PATCH /v1/catalogs/products/{product_id}
func (c *Client) UpdateProduct(productID string, patchObject []*PatchObject) error {
	req, err := c.NewRequest("PATCH", fmt.Sprintf("%s%s", c.APIBase, "/v1/catalogs/products/"+productID), patchObject)
	if err != nil {
		return err
	}

	return c.SendWithBasicAuth(req, nil)
}
