OAPIDIR = api/openapi
VOUCHER_PORTDIR = internal/voucher/port/public
.PHONY: openapi
openapi: openapi_http 

.PHONY: openapi_http
openapi_http:
		openapi-generator generate -g go-gin-server -o ${VOUCHER_PORTDIR} -i ${OAPIDIR}/voucher.yml -c ${OAPIDIR}/config.json -t ${OAPIDIR}/templates/voucher