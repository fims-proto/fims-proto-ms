package http

// TODO remove: test purpose
// TODO remove: this class cannot be used in production as it's not following our design pattern

import (
	"github/fims-proto/fims-proto-ms/internal/ledger/adapter"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type Handler struct {
	repo adapter.LedgerMemoryRepository
}

func NewHandler(repo adapter.LedgerMemoryRepository) Handler {
	return Handler{repo: repo}
}

func (h Handler) allLedgers(c *gin.Context) {
	ls, _ := h.repo.AllLedgers(c)

	res := []Ledger{}
	for _, l := range ls {
		res = append(res, Ledger{
			Sob:            l.Sob(),
			Number:         l.Number(),
			Title:          l.Title(),
			SuperiorNumber: l.SuperiorNumber(),
			AccountType:    l.AccountType().String(),
			Debit:          l.Debit(),
			Credit:         l.Credit(),
			Balance:        l.Balance(),
		})
	}

	c.JSON(http.StatusOK, res)
}

func InitRouter(h Handler, r *gin.Engine) {
	g := r.Group("/ledgers-test")
	{
		g.GET("/", h.allLedgers)
	}
}

type Ledger struct {
	Sob            string
	Number         string
	Title          string
	SuperiorNumber string
	AccountType    string
	Debit          decimal.Decimal
	Credit         decimal.Decimal
	Balance        decimal.Decimal
}
