package requests

import (
	"net/http"
)

type UpdateExpenseRequest struct {
	*CreateExpenseRequest
}

func (u *UpdateExpenseRequest) Bind(r *http.Request) error {
	return nil

	return u.CreateExpenseRequest.Bind(r)
}