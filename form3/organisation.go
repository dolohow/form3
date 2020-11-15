package form3

import "net/http"

type organisationService struct {
	Accounts *accountsService
}

func newOrganisationService(httpClient *http.Client, baseURL string) *organisationService {
	return &organisationService{
		Accounts: newAccountsSerivce(httpClient, baseURL),
	}
}
