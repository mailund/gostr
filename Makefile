tests:
	go test -v ./...

cover: 
	go test ./... -coverprofile cover.out
	go tool cover -func cover.out

.phony: tests cover
