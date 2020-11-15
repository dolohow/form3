package form3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const resourceURL = "/v1/organisation/accounts"

// Attributes represents Organization Account Attributes.
type Attributes struct {
	Country                 string   `json:"country"`
	BaseCurrency            string   `json:"base_currency"`
	AccountNumber           string   `json:"account_number"`
	BankID                  string   `json:"bank_id"`
	BankIDCode              string   `json:"bank_id_code"`
	Bic                     string   `json:"bic"`
	Iban                    string   `json:"iban"`
	Name                    []string `json:"name"`
	AlternativeNames        []string `json:"alternative_names"`
	AccountClassification   string   `json:"account_classification"`
	JointAccount            bool     `json:"joint_account"`
	AccountMatchingOptOut   bool     `json:"account_matching_opt_out"`
	SecondaryIdentification string   `json:"secondary_identification"`
	Switched                bool     `json:"switched"`
	Status                  string   `json:"status"`
}

// Account represents Organization Account.
type Account struct {
	Type           string      `json:"type"`
	ID             string      `json:"id"`
	OrganisationID string      `json:"organisation_id"`
	Version        int         `json:"version"`
	Attributes     *Attributes `json:"attributes"`
}

type topLevelJSONProperty struct {
	Data *Account `json:"data"`
}

type topLevelJSONArray struct {
	Data []*Account `json:"data"`
}

type accountsService struct {
	httpClient *http.Client
	baseURL    string
}

func extractData(body []byte) *Account {
	data := &topLevelJSONProperty{}
	json.Unmarshal(body, data)
	return data.Data
}

func newAccountsSerivce(httpClient *http.Client, baseURL string) *accountsService {
	return &accountsService{httpClient: httpClient, baseURL: baseURL}
}

func (s *accountsService) getResourceURL() string {
	return s.baseURL + resourceURL
}

// List fetches all Organisation Accounts.
func (s *accountsService) List(params *URLParameters) ([]*Account, error) {
	url := s.getResourceURL() + encodeURLParameters((params))
	res, err := s.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	data := &topLevelJSONArray{}
	json.Unmarshal(body, data)

	return data.Data, checkAPIError(res, body)
}

// Fetch fetches a single Account from Organisation based on passed id.
func (s *accountsService) Fetch(id string) (*Account, error) {
	url := fmt.Sprintf("%s/%s", s.getResourceURL(), id)
	res, err := s.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return extractData(body), checkAPIError(res, body)
}

// Create creates new Account in Organisation.
func (s *accountsService) Create(account *Account) (*Account, error) {
	data := &topLevelJSONProperty{Data: account}
	json, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	res, err := s.httpClient.Post(s.getResourceURL(), "application/vnd.api+json", bytes.NewBuffer(json))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return extractData(body), checkAPIError(res, body)
}

// Delete removed a single version of Organisation Account resource.
func (s *accountsService) Delete(id string, version int) error {
	query := url.Values{}
	query.Set("version", strconv.Itoa(version))
	url := fmt.Sprintf("%s/%s?%s", s.getResourceURL(), id, query.Encode())

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	res, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return checkAPIError(res, body)
}
