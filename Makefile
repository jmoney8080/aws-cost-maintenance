all: build

dependencies: go.mod
	go get -u ./...

verify: dependencies
	go mod verify

build: verify main.go
	GOOS="linux" GOARCH="amd64" go build -o build/aws-cost-maintenance_linux-amd64
	GOOS="darwin" GOARCH="amd64" go build -o build/aws-cost-maintenance_darwin-amd64

clean:
	rm -rvf build/*
