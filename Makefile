.PHONY: test
test:
	go test ./... -count=1

.PHONY: peg
peg:
	peg -output internal/common/data/filterable/filterable_ast.go internal/common/data/filterable/filterable.peg

.PHONY: swag
swag:
	swag init -g api/api.go -o docs/swagger_generated --parseDependencyLevel 1

.PHONY: fmt
fmt:
	swag fmt
	gofumpt -l -w internal/ cmd/

.PHONY: lint
lint:
	golangci-lint run ./...