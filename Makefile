tests:
	go test ./...

lint:
	golangci-lint run

cover: 
	go test ./... -coverprofile cover.out
	go tool cover -func cover.out

.phony: tests cover
