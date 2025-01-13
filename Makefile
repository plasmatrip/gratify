
.PHONY : all
all: preproc build-all

.PHONY : preproc
preproc: clean fmt test check-coverage 

.PHONY : build-all
build-all: clean gophermart run-autotests

.PHONY : autotests
autotests: run-autotests

gophermart:
	go build -o ./cmd/gophermart/gophermart ./cmd/gophermart/main.go

test:
	go test ./... -race -coverprofile=cover.out -covermode=atomic

.PHONY : clean
clean:
	-rm ./cmd/gophermart/gophermart 2>/dev/null
	-rm ./cover.out 2>/dev/null
	-rm ./backup.dat 2>dev/null

check-coverage:
	go tool cover -html cover.out -o cover.html

.PHONY : fmt
fmt:
	go fmt ./...
	goimports -v -w .

# .PHONY : lint
# lint:
# 	golangci-lint run ./...

RUN_ADDRESS := 8080

.PHONY : run-autotestsg
# run-autotests: iter13
run-autotests: all_test


.PHONY : all_test
all_test:
	 gophermarttest \
            -test.v -test.run=^TestGophermart$$ \
            -gophermart-binary-path=cmd/gophermart/gophermart \
            -gophermart-host=localhost \
            -gophermart-port=8080 \
            -gophermart-database-uri="postgres://gratify:password@localhost:5432/gratify?sslmode=disable" \
            -accrual-binary-path=cmd/accrual/accrual_linux_amd64 \
            -accrual-host=localhost \
            -accrual-port=8082 \
            -accrual-database-uri="postgresql://postgres:postgres@postgres/praktikum?sslmode=disable"