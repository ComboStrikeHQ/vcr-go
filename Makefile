all: lint test

deps:
	go get -u github.com/stretchr/testify github.com/alecthomas/gometalinter
	gometalinter -iu

lint:
	gometalinter -D errcheck --deadline=15s

test:
	go test -v -race -cover |tee testlog.out
	grep "coverage: 100.0% of statements" testlog.out
	@rm testlog.out
