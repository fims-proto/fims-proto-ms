package query

type Account struct {
	Sob             string
	Number          string
	Title           string
	AccountType     string
	SuperiorAccount *Account
}
