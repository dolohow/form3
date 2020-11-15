package form3

import (
	"net/http"
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("%s not equal to %s", a, b)
	}
}

func assertNotEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("%s is equal to %s", a, b)
	}
}

func assertNil(t *testing.T, a interface{}) {
	if a != nil {
		t.Errorf("%s is not nil", a)
	}
}

func assertTrue(t *testing.T, a bool, message string) {
	if !a {
		t.Errorf(message, a)
	}
}

func assertMethod(t *testing.T, req *http.Request, method string) {
	assertEqual(t, req.Method, method)
}

func assertURL(t *testing.T, req *http.Request, url string) {
	assertEqual(t, req.URL.String(), url)
}
