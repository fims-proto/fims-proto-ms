package query

type Account struct {
	Number          string
	Title           string
	AccountType     string
	SuperiorAccount *Account
}
