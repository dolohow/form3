package form3

import "testing"

func TestForm3_NewClient(t *testing.T) {
	client := NewClient("http://localhost:8080")
	assertNotEqual(t, client, nil)
}
