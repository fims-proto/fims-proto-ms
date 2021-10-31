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
- kratos, use `http://127.0.0.1:4455/kratos/public/<any path>`
- fims, use `http://127.0.0.1:4455/fims/s/<any path>`

## Create user
In [fims-iap](https://github.com/fims-proto/fims-iap) folder, run shell `./scripts/user_creation.sh`, then `./scripts/user_invitation.sh`
