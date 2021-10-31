# fims-proto-ms
## Local dev steps:
1. in fims-proto-ms folder, start FIMS backend ms + postgres with `docker compose up -build`
2. in [fims-iap](https://github.com/fims-proto/fims-iap) folder, start ory/kratos + ory/oathkeeper with `docker compose up`
3. in [fims-proto-ui](https://github.com/fims-proto/fims-proto-ui) folder, start vite in 'dev' mode with `npm run dev`
4. ready to go :)

## Frontend
Port `5000` is used in local dev environment: `http://127.0.0.1:5000`.  
In production mode, frontend should also be wrapped into fims-iap, to avoid common cors issue.  

## Backend
Port `4455` is exposed by ory/oathkeeper as reverse proxy, acting as single entry point.  
To access:
- kratos public API, use `http://127.0.0.1:4455/kratos/public/<any path>`
- fims public API, use `http://127.0.0.1:4455/fims/s/<any path>`

## Create user
In [fims-iap](https://github.com/fims-proto/fims-iap) folder, run shell `./scripts/user_creation.sh`, then `./scripts/user_invitation.sh`

## Swagger UI
After starting fims-proto-ms + fims-iap + fims-proto-ui, login fims-proto-ui, then visit `http://127.0.0.1:4455/fims/public/swagger.index`.  

## Postman test
1. Import [postman collection](https://github.com/fims-proto/fims-proto-ms/tree/master/pm_collection) into Postman app
2. Get jwt token by visit `http://127.0.0.1:5000/devops/jwt`, login if necessary, copy all in the page
3. In postman collection, choose Authorization type as `Bearer Token`, put the jwt token into token field.
