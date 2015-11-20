package paypalsdk

import (
	"testing"
)

func TestGetAccessToken(t *testing.T) {
	c, _ := NewClient("clid", "secret", APIBaseSandBox)
	c.GetAccessToken()
}
