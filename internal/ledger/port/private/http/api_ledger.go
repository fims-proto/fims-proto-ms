package http

// TODO remove: test purpose
// TODO remove: this class cannot be used in production as it's not following our design pattern

import (
	"github/fims-proto/fims-proto-ms/internal/ledger/adapter"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"
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

func (h Handler) loadLedgers(c *gin.Context) {
	l1, _ := domain.NewLedger("10000101", "1000-01-01", "100001", 1)
	l2, _ := domain.NewLedger("100001", "1000-01", "1000", 1)
	l3, _ := domain.NewLedger("1000", "1000", "", 1)
	l4, _ := domain.NewLedger("20000202", "2000-02-02", "200002", 3)
	l5, _ := domain.NewLedger("200002", "2000-02", "2000", 3)
	l6, _ := domain.NewLedger("2000", "2000", "", 3)
	ls := []*domain.Ledger{l1, l2, l3, l4, l5, l6}

	for _, l := range ls {
		if err := h.repo.AddLedger(c.Request.Context(), l); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (h Handler) allLedgers(c *gin.Context) {
	ls, _ := h.repo.AllLedgers(c.Request.Context())

	res := []Ledger{}
	for _, l := range ls {
		res = append(res, Ledger{
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
		g.POST("/load", h.loadLedgers)
	}
}

type Ledger struct {
	Number         string
	Title          string
	SuperiorNumber string
	AccountType    string
	Debit          decimal.Decimal
	Credit         decimal.Decimal
	Balance        decimal.Decimal
}
