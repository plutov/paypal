package paypal

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// testClientID, testSecret imported from order_test.go

// All test values are defined here
// var testClientID = "AXy9orp-CDaHhBZ9C78QHW2BKZpACgroqo85_NIOa9mIfJ9QnSVKzY-X_rivR_fTUUr6aLjcJsj6sDur"
// var testSecret = "EBoIiUSkCKeSk49hHSgTem1qnjzzJgRQHDEHvGpzlLEf_nIoJd91xu8rPOBDCdR_UYNKVxJE-UgS2iCw"
var testUserID = "https://www.paypal.com/webapps/auth/identity/user/VBqgHcgZwb1PBs69ybjjXfIW86_Hr93aBvF_Rgbh2II"
var testCardID = "CARD-54E6956910402550WKGRL6EA"

var testProductId = ""   // will be fetched in  func TestProduct(t *testing.T)
var testBillingPlan = "" // will be fetched in  func TestSubscriptionPlans(t *testing.T)

const alphabet = "abcedfghijklmnopqrstuvwxyz"

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
	//create a product
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
	// c.Lock()
	c.mu.Lock()
	// Note: Here we do not want to `defer c.Unlock()` because we need `c.Send(...)`
	// to happen outside of the locked section.

	if c.mu.TryLock() {
		// if the code is able to acquire a lock
		// despite the mutex of c being locked, throw an error
		err := errors.New("TryLock succeeded inside sendWithAuth with mutex locked")
		return err
	}

	if c.Token != nil {
		if !c.tokenExpiresAt.IsZero() && c.tokenExpiresAt.Sub(time.Now()) < RequestNewTokenBeforeExpiresIn {
			// c.Token will be updated in GetAccessToken call
			if _, err := c.GetAccessToken(req.Context()); err != nil {
				// c.Unlock()
				c.mu.Unlock()
				return err
			}
		}

		req.Header.Set("Authorization", "Bearer "+c.Token.Token)
	}
	// Unlock the client mutex before sending the request, this allows multiple requests
	// to be in progress at the same time.
	// c.Unlock()
	c.mu.Unlock()

	if !c.mu.TryLock() {
		// if the code is unable to acquire a lock
		// despite the mutex of c being unlocked, throw an error
		err := errors.New("TryLock failed inside sendWithAuth with mutex unlocked")
		return err
	}
	c.mu.Unlock() // undo changes from the previous TryLock

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
	c.GetAccessToken(context.Background())

	// Basic Testing of the private mutex field
	c.mu.Lock()
	if c.mu.TryLock() {
		t.Fatalf("TryLock succeeded with mutex locked")
	}
	c.mu.Unlock()
	if !c.mu.TryLock() {
		t.Fatalf("TryLock failed with mutex unlocked")
	}
	c.mu.Unlock() // undo changes from the previous TryLock

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
