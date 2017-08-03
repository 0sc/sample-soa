export DATABASE_URL = root:abc123@tcp(127.0.0.1:3306)/USRINFO
export BOOK_SVC_ADDR = localhost:8899
export BOOK_PRICE_SVC_ADDR = localhost:9999
export PORT = :8999

all:
	go run *.go

gen_proto:
	cd ../book-svc/rpc && \
		protoc --go_out=plugins=grpc:../../auth-svc/rpc/book book.proto
	cd ../book-price-svc/rpc && \
		protoc --go_out=plugins=grpc:../../auth-svc/rpc/book-price book-price.proto
		

		