GO ?= go

build:
	$(GO) build -v 

test:
	$(GO) test -v -cover .

clean:
	$(GO) clean
