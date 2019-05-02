.PHONY: test

test:
	go test -v -race -count=1 ./test/...
