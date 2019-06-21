package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go-rest-api/src/github.com/sethpawas/go-rest-api/crud_Interface"
	"go-rest-api/src/github.com/sethpawas/go-rest-api/errrs"
	"go-rest-api/src/github.com/sethpawas/go-rest-api/requests"
	"go-rest-api/src/github.com/sethpawas/go-rest-api/responses"
	"go-rest-api/src/github.com/sethpawas/go-rest-api/types"
	"log"
	"net/http"
	"time"
)

var expenses types.Expenses
var expense types.Expense
var temp types.Expense
var err error


type Mysql struct{
	Db *gorm.DB
}


func main() {

	db, err := gorm.Open("mysql", "root:root@tcp(localhost:3306)/expense?charset=utf8&parseTime=True")

	if err != nil {
		fmt.Println("Connection not established")
	}

	if (!db.HasTable(&types.Expense{})) { //if no db then create db
		db.AutoMigrate(&types.Expense{})
	}

	sel := &Mysql{db}//for interface connection

	Init(sel)
}


func Init(d crud_Interface.Database){
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/expenses", func(r chi.Router) {
		r.Post("/", d.CreateExpense)
		r.Get("/", d.ListAllExpense)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(d.ArticleCtx)
			r.Get("/", d.GetId)
			r.Put("/", d.UpdateExp)
			r.Delete("/", d.DeleteExp)
		})
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}


func (db *Mysql) ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		expenseID := chi.URLParam(r, "id")
		//Db, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/expense?charset=utf8&parseTime=True")
		db:= db.Db.Table("expenses").Where("id = ?", expenseID).Find(&expense)

		ctx := context.WithValue(r.Context(), "expense", db )
		next.ServeHTTP(w, r.WithContext(ctx))

		if db.RowsAffected==0{
			err=errors.New("ID not Found")
			render.Render(w, r, errrs.ErrRender(err))
		}
	})
}


func (db *Mysql) CreateExpense(writer http.ResponseWriter, request *http.Request) {

	var req requests.CreateExpenseRequest

	//Db, err := gorm.Open("mysql", "root:root@tcp(localhost:3306)/expense?charset=utf8&parseTime=True")
	//defer Db.Close()

	err = render.Bind(request, &req)
	if err != nil {
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}

	temp:=*req.Expense
	temp.CreatedOn=time.Now()
	temp.UpdatedOn=time.Now()


	db.Db.Create(&temp)

	_,_ = fmt.Fprintln(writer, `{"success": true}`)
	//render.Render(writer, request, response.Getoneresponse(req.Expense))
	//expenses = append(expenses, *req.Expense)
	render.Render(writer, request, responses.List1expense(req.Expense))
	//db.Db.Close()
}


func (db *Mysql) UpdateExp (writer http.ResponseWriter, request *http.Request) {

	db.Db =request.Context().Value("expense").(*gorm.DB)

	var req requests.UpdateExpenseRequest

	err= render.Bind(request,&req)
	if err != nil {
		panic("Error occurred")
	}

	var temp types.Expense
	temp=*req.Expense
	temp.UpdatedOn=time.Now()

	Db := db.Db.Update(&temp)
	if(Db.RowsAffected == 0){
		err:=errors.New("Expense not found")
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}else{
		errs:=render.Render(writer, request, responses.List1expense(&temp))
		if errs != nil {
		render.Render(writer,request,errrs.ErrRender(errs))
			return
		}

	}
}


func (db *Mysql) GetId(writer http.ResponseWriter, request *http.Request) {

	db.Db =request.Context().Value("expense").(*gorm.DB)

	if err != nil {
		panic("Error occurred")
	}

	var temp types.Expense
	Db:= db.Db.Find(&temp)
	if(Db.RowsAffected == 0){
		err:=errors.New("Expense not found")
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}else{
		render.Render(writer, request, responses.List1expense(&temp))
		return
	}
}


func (db *Mysql) ListAllExpense(writer http.ResponseWriter, request *http.Request) {

	//db, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/expense?charset=utf8&parseTime=True")
	//defer Db.Close()
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Connection established")
	}

	var temp types.Expenses
	db.Db.Table("expenses").Find(&temp)
	err = render.Render(writer, request, responses.NewExpensesResponse(&temp))
	if err != nil {
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}
	db.Db.Close()
}


func (db *Mysql) DeleteExp(writer http.ResponseWriter, request *http.Request) {

	db.Db =request.Context().Value("expense").(*gorm.DB)//not gorm.DB because it takes the pointer. it takes the value of from the Articlectx
	// and if the id doesnt exist it handles id if doesnt exist and put the id to expense

	if err != nil {
		panic("Error occurred")
	}

	Db:= db.Db.Delete(&temp)
	if(Db.RowsAffected == 0){
		err:=errors.New("Expense not found")
		render.Render(writer,request,errrs.ErrRender(err))
		return
	}else{
		fmt.Fprintf(writer,"sucessful delete")
		return
	}
}