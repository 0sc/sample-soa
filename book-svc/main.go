package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/proj/book-svc/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var db *sql.DB
var err error

type rpcServer struct{}

func (*rpcServer) GetBooks(ctx context.Context, in *rpc.Empty) (*rpc.Books, error) {
	defer db.Close()

	st, err := db.Prepare("SELECT Title, Author , ISBN , Description FROM books")
	if err != nil {
		fmt.Print(err)
	}
	if err != nil {
		return nil, err
	}
	rows, err := st.Query()
	if err != nil {
		fmt.Print(err)
	}

	var result = &rpc.Books{}
	var book rpc.Book
	var title, author, isbn, description string

	for rows.Next() {
		err = rows.Scan(&title, &author, &isbn, &description)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		book = rpc.Book{Title: title, Author: author, Isbn: isbn, Description: description}
		result.Books = append(result.Books, &book)
	}

	return result, nil
}

func setupRPCServer() *grpc.Server {
	grpcServer := grpc.NewServer()
	rpcServer := &rpcServer{}
	rpc.RegisterBookSvcServer(grpcServer, rpcServer)

	return grpcServer
}

func main() {
	db, err = sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	server := setupRPCServer()
	port := os.Getenv("PORT")
	listener, err := net.Listen("tcp", port)

	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
	fmt.Println("Book service is running on: ", port)

	server.Serve(listener)
}
