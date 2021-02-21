.PHONY: openapi
openapi: openapi_http 

.PHONY: openapi_http
openapi_http:
	oapi-codegen -generate types -o internal/voucher/port/public/http/openapi_types.gen.go -package http api/openapi/voucher.yml

.PHONY: fmt
fmt:
	gofumpt -l -w internal/ cmd/ 