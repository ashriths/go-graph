.PHONY: build all rall fmt tags lc doc


all: build
	go install ./src/...

rall:
	go build -a ./src/...

fmt:
	gofmt -s -w -l .

tags:
	gotags `find . -name "*.go"` > tags

lc:
	wc -l `find . -name "*.go"`

doc:
	godoc -http=:8000

build:
	go get github.com/google/uuid

