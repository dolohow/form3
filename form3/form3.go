package form3

import (
	"fmt"
	"net/http"
	"reflect"
)

// Client is a client for making API requests.
type Client struct {
	URL          string
	Organisation *organisationService
}

// URLParameters holds common query params that could be passed to API endpoint.
type URLParameters struct {
	PageNumber string `url:"page[number]"`
	PageSize   string `url:"page[size]"`
}

func encodeURLParameters(params *URLParameters) (query string) {
	if params == nil {
		return ""
	}

	t := reflect.TypeOf(*params)
	v := reflect.ValueOf(params)

	query += "?"

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("url")
		value := v.Elem().Field(i).String()
		if i+1 == t.NumField() {
			query = fmt.Sprintf("%s%s=%s", query, tag, value)
		} else {
			query = fmt.Sprintf("%s%s=%s&", query, tag, value)
		}
	}

	return
}

// NewClient returns new Client.
func NewClient(url string) *Client {
	httpClient := &http.Client{}
	return &Client{
		URL:          url,
		Organisation: newOrganisationService(httpClient, url),
	}
}
