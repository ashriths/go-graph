.PHONY: build all rall fmt tags lc doc


all: build
	go install ./...

rall:
	go build -a ./...

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
	go get github.com/samuel/go-zookeeper/zk

test:
	go test ./...

testv:
	go test -v ./...
