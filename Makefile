.PHONY: build
build:
	go test .
	go vet
	go build .

.PHONY:
update:
	go get -u
	go mod tidy
