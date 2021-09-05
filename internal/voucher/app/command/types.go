package command

import "github.com/google/uuid"

type LineItemCmd struct {
	Id            uuid.UUID
	Summary       string
	AccountNumber string
	Debit         string
	Credit        string
}
