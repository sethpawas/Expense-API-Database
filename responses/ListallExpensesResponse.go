package responses

import (
	"go-rest-api/src/github.com/sethpawas/go-rest-api/types"
	"net/http"
)

type ExpensesResponse struct {
  Expenses *types.Expenses
}

func NewExpensesResponse(expenses *types.Expenses) *ExpensesResponse{
	return &ExpensesResponse{Expenses: expenses}

}

func (e *ExpensesResponse) Render(w http.ResponseWriter, r *http.Request) error {

	return nil
}