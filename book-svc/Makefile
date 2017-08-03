export DATABASE_URL = root:abc123@tcp(127.0.0.1:3306)/BOOKINFO
export PORT = :8899

all:
	go run main.go

gen_proto:
	cd rpc && protoc --go_out=plugins=grpc:. book.proto
