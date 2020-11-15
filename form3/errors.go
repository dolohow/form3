package form3

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError Represents error returned from API.
type APIError struct {
	ErrorMessage string `json:"error_message"`
	StatusCode   int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api: %d: %v", e.StatusCode, e.ErrorMessage)
}

func checkAPIError(res *http.Response, body []byte) error {
	if res.StatusCode >= 400 {
		var apiError APIError
		json.Unmarshal(body, &apiError)
		apiError.StatusCode = res.StatusCode
		return &apiError
	}
	return nil
}
