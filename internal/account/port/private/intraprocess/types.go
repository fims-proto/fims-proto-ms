package intraprocess

type AccountResponse struct {
	Number          string
	Title           string
	AccountType     string
	SuperiorAccount *AccountResponse
}
