OAPIDIR = api/openapi
VOUCHER_PORTDIR = internal/voucher/port/public
.PHONY: openapi
openapi: openapi_http 

.PHONY: openapi_http
openapi_http:
		openapi-generator generate --global-property=models,apis -g go-server -o ${VOUCHER_PORTDIR} -i ${OAPIDIR}/voucher.yml -c api/openapi/config.json