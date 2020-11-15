package form3

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
)

const fixturesDir = "fixtures/Organisation/Accounts"

func getFakeCreateAccountData() *Account {
	return &Account{
		ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
		Type:           "accounts",
		OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
		Attributes: &Attributes{
			Country:               "GB",
			AccountClassification: "Personal",
		},
	}
}

func createServer(c func(http.ResponseWriter, *http.Request, []byte), fixtureName string) (*accountsService, *httptest.Server) {

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if fixtureName != "" {
			data, _ := ioutil.ReadFile(path.Join(fixturesDir, fixtureName))
			data = bytes.ReplaceAll(data, []byte(" "), []byte(""))
			data = bytes.ReplaceAll(data, []byte("\n"), []byte(""))
			w.Write(data)
			c(w, r, data)
		} else {
			c(w, r, nil)
		}
	})
	server := httptest.NewServer(handler)

	return newAccountsSerivce(server.Client(), server.URL), server
}

func TestAccount_List(t *testing.T) {
	accountService, server := createServer(func(w http.ResponseWriter, r *http.Request, data []byte) {
		assertMethod(t, r, "GET")
		assertURL(t, r, "/v1/organisation/accounts")
	}, "request_list.json")
	defer server.Close()

	ret, err := accountService.List(nil)

	assertNil(t, err)
	assertEqual(t, ret[0].ID, "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc")
	assertEqual(t, ret[0].OrganisationID, "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c")
}

func TestAccount_ListParameters(t *testing.T) {
	accountService, server := createServer(func(w http.ResponseWriter, r *http.Request, data []byte) {
		assertMethod(t, r, "GET")
		assertURL(t, r, "/v1/organisation/accounts?page[number]=2&page[size]=10")
	}, "request_list.json")
	defer server.Close()

	ret, err := accountService.List(&URLParameters{PageNumber: "2", PageSize: "10"})

	assertNil(t, err)
	assertEqual(t, ret[0].ID, "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc")
	assertEqual(t, ret[0].OrganisationID, "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c")
}

func TestAccount_ListError(t *testing.T) {
	accountService, server := createServer(func(w http.ResponseWriter, r *http.Request, data []byte) {
		w.WriteHeader(500)
	}, "")
	defer server.Close()

	_, err := accountService.List(nil)

	assertEqual(t, err.(*APIError).StatusCode, 500)
}

func TestAccount_Fetch(t *testing.T) {
	accountService, server := createServer(func(w http.ResponseWriter, r *http.Request, data []byte) {
		assertMethod(t, r, "GET")
		assertURL(t, r, "/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc")
	}, "request_fetch.json")
	defer server.Close()

	ret, err := accountService.Fetch("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc")

	assertNil(t, err)
	assertEqual(t, ret.ID, "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc")
	assertEqual(t, ret.OrganisationID, "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c")
}

func TestAccount_FetchError(t *testing.T) {
	accountService, server := createServer(func(w http.ResponseWriter, r *http.Request, data []byte) {
		w.WriteHeader(500)
	}, "")
	defer server.Close()

	_, err := accountService.Fetch("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc")

	assertEqual(t, err.(*APIError).StatusCode, 500)
}

func TestAccount_Create(t *testing.T) {
	accountService, server := createServer(func(w http.ResponseWriter, r *http.Request, data []byte) {
		body, _ := ioutil.ReadAll(r.Body)
		assertTrue(t, bytes.Equal(data, body), "Request does not match")
		assertMethod(t, r, "POST")
		assertURL(t, r, "/v1/organisation/accounts")
	}, "request_create.json")
	defer server.Close()

	ret, err := accountService.Create(getFakeCreateAccountData())

	assertNil(t, err)
	assertEqual(t, ret.ID, "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc")
	assertEqual(t, ret.OrganisationID, "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c")
}

func TestAccount_CreateError(t *testing.T) {
	accountService, server := createServer(func(w http.ResponseWriter, r *http.Request, data []byte) {
		w.WriteHeader(409)
		w.Write([]byte(`{"error_message":"Resource exists"}`))
	}, "")
	defer server.Close()

	_, err := accountService.Create(getFakeCreateAccountData())

	assertEqual(t, err.(*APIError).StatusCode, 409)
	assertEqual(t, err.(*APIError).ErrorMessage, "Resource exists")
}

func TestAccount_Delete(t *testing.T) {
	accountService, server := createServer(func(w http.ResponseWriter, r *http.Request, data []byte) {
		w.WriteHeader(204)
		assertMethod(t, r, "DELETE")
		assertURL(t, r, "/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc?version=5")
	}, "")
	defer server.Close()

	err := accountService.Delete("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc", 5)

	assertNil(t, err)
}

func TestAccount_DeleteError(t *testing.T) {
	accountService, server := createServer(func(w http.ResponseWriter, r *http.Request, data []byte) {
		w.WriteHeader(404)
	}, "")
	defer server.Close()

	err := accountService.Delete("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc", 5)

	assertEqual(t, err.(*APIError).StatusCode, 404)
}
