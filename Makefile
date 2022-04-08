.PHONY: install
install:
	go install github.com/alvaroloes/enumer@v1.1.2
	go install github.com/go-bindata/go-bindata/go-bindata@v3.1.2+incompatible
	go install github.com/golang/protobuf/protoc-gen-go@v1.1.0
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.28.1
	go install github.com/vektra/mockery/cmd/mockery@v1.0.0
	go install github.com/volatiletech/sqlboiler@v3.7.1+incompatible
	go install github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql@v3.7.1+incompatible
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
	go get -t


.PHONE: clean
clean:
	go fmt ./...
	go mod tidy


.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/pokemon.proto


.PHONY: run
run: proto
	go run src/main.go


.PHONY: models
models:
	cd database; sqlboiler --wipe mysql; cd ..;