.PHONY: swag
swag:
	swag init -g api/api.go

.PHONY: fmt
fmt:
	gofumpt -l -w internal/ cmd/

.PHONY: lint
lint:
	golangci-lint run