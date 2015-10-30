package paypalsdk

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	_, err := NewClient("", "", "")
	if err == nil {
		t.Errorf("All arguments are required in NewClient()")
	}
}
