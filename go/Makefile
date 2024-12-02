isuports: go.mod go.sum *.go cmd/isuports/*
	go build -o isuports ./cmd/isuports

dbclean: go.mod go.sum *.go cmd/dbclean/*
	go build -o cleandb ./cmd/dbclean

test:
	go test -v ./...
