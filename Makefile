all: lint test

deps:
	go get -u github.com/stretchr/testify github.com/alecthomas/gometalinter
	gometalinter -iu

lint:
	gometalinter -E gofmt -D errcheck -D gosec

test:
	go test -v -race -cover |tee testlog.out
	grep "coverage: 100.0% of statements" testlog.out
	@rm testlog.out
