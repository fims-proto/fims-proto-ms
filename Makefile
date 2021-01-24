.PHONY: openapi
openapi: openapi_http 

.PHONY: openapi_http
openapi_http:
		openapi-codegen 