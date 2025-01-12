
.PHONY : all
all: preproc build-all

.PHONY : preproc
preproc: clean fmt test check-coverage 

.PHONY : build-all
build-all: clean gophermart run-autotests

.PHONY : autotests
autotests: run-autotests

gophermart:
	go build -o ./bin/gophermart ./cmd/gophermart/main.go

test:
	go test ./... -race -coverprofile=cover.out -covermode=atomic

.PHONY : clean
clean:
	-rm ./bin/gophermart 2>/dev/null
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
run-autotests: all_test race-condition


.PHONY : all_test
all_test:
	 metricstest -test.run=^TestIteration14$$ \
			-agent-binary-path=./bin/agent \
			-binary-path=./bin/gophermart \
			-database-dsn='postgres://metriq:password@localhost:5432/metriq?sslmode=disable' \
			 -key="${TEMP_FILE}" \
			-gophermart-port=$(gophermart_PORT) \
			-source-path=.

.PHONY : race-condition
race-condition:
	go test -v -race ./...