package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/proj/book-price-svc/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var db *sql.DB
var err error

type rpcServer struct{}

func (s *rpcServer) GetPrice(ctx context.Context, payload *rpc.Payload) (*rpc.Payload, error) {
	bookname := payload.GetValue()

	var databaseBookprice string

	err := db.QueryRow("SELECT price FROM bookprice WHERE bookname=?", bookname).Scan(&databaseBookprice)

	if err != nil {
		return nil, err
	}

	return &rpc.Payload{Value: databaseBookprice}, nil
}

func setupRPCServer() *grpc.Server {
	grpcServer := grpc.NewServer()
	rpcServer := &rpcServer{}
	rpc.RegisterBookPriceSvcServer(grpcServer, rpcServer)

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
	fmt.Println("Book Price service is running on: ", port)

	server.Serve(listener)
}
