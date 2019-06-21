package crud_Interface

import "net/http"

type Database interface{
	ArticleCtx(next http.Handler) http.Handler
	CreateExpense(writer http.ResponseWriter, request *http.Request)
	UpdateExp(writer http.ResponseWriter, request *http.Request)
	GetId(writer http.ResponseWriter, request *http.Request)
	ListAllExpense(writer http.ResponseWriter, request *http.Request)
	DeleteExp(writer http.ResponseWriter, request *http.Request)
}

