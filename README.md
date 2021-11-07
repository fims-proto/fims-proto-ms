# fims-proto-ms
## Local dev steps:
1. in fims-proto-ms folder, start postgres with `docker compose up -build`
2. in fims-proto-ms folder, start fims backend service with `go run cmd/main.go`, running at port `5002`
3. in [fims-iap-dev](https://github.com/fims-proto/fims-iap-dev) folder, start ory/oathkeeper with `docker compose up -build`, running at port `4455`
4. in [fims-proto-ui](https://github.com/fims-proto/fims-proto-ui) folder, start vite in 'dev' mode with `npm run dev`, running at port `5001`
5. Open [`http://127.0.0.1:4455/ui/`](http://127.0.0.1:4455/ui/), ready to go :)

为了开发环境的简单便捷，目前 local dev 环境暂不引入权鉴 (ory/kratos).   

## Frontend
Port `5001` is used in local dev environment: `http://127.0.0.1:5001`.  

## Backend
Port `4455` is exposed by ory/oathkeeper as reverse proxy, acting as single entry point.  
To access:
- fims backend API, use `http://127.0.0.1:4455/fims/**`
- fims frontend, use `http://127.0.0.1:4455/ui/**`

## Swagger UI
After starting fims-proto-ms, visit [`http://127.0.0.1:4455/fims/swagger/index.html`](http://127.0.0.1:4455/fims/swagger/index.html).  

## Postman test
Import [postman collection](https://github.com/fims-proto/fims-proto-ms/tree/master/pm_collection) into Postman app
