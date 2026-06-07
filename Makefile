.PHONY: peg
peg:
	peg -output internal/common/data/filterable/filterable_ast.go internal/common/data/filterable/filterable.peg

.PHONY: test
test:
	go test ./... -count=1

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: swag
swag:
	go run cmd/main.go --gen-spec

.PHONY: fmt
fmt:
	gofumpt -l -w internal/ cmd/
