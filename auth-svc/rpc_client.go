package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/context"

	book "github.com/proj/auth-svc/rpc/book"
	price "github.com/proj/auth-svc/rpc/book-price"
	"google.golang.org/grpc"
)

var booksRPC book.BookSvcClient
var priceRPC price.BookPriceSvcClient

func setupRPCClientConnection() {
	booksRPC = book.NewBookSvcClient(rpcConnectionToBookSvc())
	priceRPC = price.NewBookPriceSvcClient(rpcConnectionToBookPriceSvc())
}

func rpcConnectionToBookSvc() *grpc.ClientConn {
	addr := os.Getenv("BOOK_SVC_ADDR")
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Couldn't connect to Book service: ", err)
	}
	return conn
}

func rpcConnectionToBookPriceSvc() *grpc.ClientConn {
	addr := os.Getenv("BOOK_PRICE_SVC_ADDR")
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Couldn't connect to Book Price service: ", err)
	}
	return conn
}

func getBooks(res http.ResponseWriter, req *http.Request) {
	books, err := booksRPC.GetBooks(context.Background(), &book.Empty{})
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(res, "An error: "+err.Error()+" occurred")
		return
	}

	res.Header().Add("Content-Type", "application/json; charset=utf-8")
	json, err := json.Marshal(books)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(res, "An error: "+err.Error()+" occurred")
		return
	}

	fmt.Fprintf(res, "%v", string(json))
}

func getBookPrice(res http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		http.Redirect(res, req, "price.html", 301)
		return
	}

	bookname := req.FormValue("bookname")
	p, err := priceRPC.GetPrice(context.Background(), &price.Payload{Value: bookname})
	if err != nil {
		fmt.Println(err)
		http.Redirect(res, req, "/price", 301)
		return
	}

	res.Write([]byte(bookname + "Price is :" + p.GetValue()))

}
