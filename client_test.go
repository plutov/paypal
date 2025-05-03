package paypal

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const alphabet = "abcedfghijklmnopqrstuvwxyz"

var (
	testClientID = "AXy9orp-CDaHhBZ9C78QHW2BKZpACgroqo85_NIOa9mIfJ9QnSVKzY-X_rivR_fTUUr6aLjcJsj6sDur"
	testSecret   = "EBoIiUSkCKeSk49hHSgTem1qnjzzJgRQHDEHvGpzlLEf_nIoJd91xu8rPOBDCdR_UYNKVxJE-UgS2iCw"
)

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func createRandomProduct(t *testing.T) Product {
	// create a product
	productData := Product{
		Name:        RandomString(10),
		Description: RandomString(100),
		Category:    ProductCategorySoftware,
		Type:        ProductTypeService,
		ImageUrl:    "https://example.com/image.png",
		HomeUrl:     "https://example.com",
	}
	return productData
}

// this is a simple copy of the SendWithAuth method, used to
// test the Lock and Unlock methods of the private mutex field
// of Client structure.
func (c *Client) sendWithAuth(req *http.Request, v interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Token == nil || (!c.tokenExpiresAt.IsZero() && time.Until(c.tokenExpiresAt) < RequestNewTokenBeforeExpiresIn) {
		if _, err := c.GetAccessToken(req.Context()); err != nil {
			return err
		}
	}
	req.Header.Set("Authorization", "Bearer "+c.Token.Token)

	return c.Send(req, v)
}

// this method is used to invoke the sendWithAuth method, which will then check
// operationally the privated mutex field of Client structure.
func (c *Client) createProduct(ctx context.Context, product Product) (*CreateProductResponse, error) {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s", c.APIBase, "/v1/catalogs/products"), product)
	response := &CreateProductResponse{}
	if err != nil {
		return response, err
	}
	err = c.sendWithAuth(req, response)
	return response, err
}

func TestClientMutex(t *testing.T) {
	c, _ := NewClient(testClientID, testSecret, APIBaseSandBox)
	_, err := c.GetAccessToken(context.Background())
	assert.NoError(t, err)

	// Operational testing of the private mutex field
	n_iter := 2

	errs := make(chan error)

	for i := 0; i < n_iter; i++ {
		go func() {
			_, err := c.createProduct(context.Background(), createRandomProduct(t))
			errs <- err
		}()
	}

	for i := 0; i < n_iter; i++ {
		err := <-errs
		assert.Equal(t, nil, err)
	}
}
