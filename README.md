# fims-proto-ms
## Local dev steps:
1. in [fims-dev-environment](https://github.com/fims-proto/fims-dev-environment) folder, start infra services (postgres, ory kratos, ory aothkeeper) with `docker compose up -build`
2. in fims-proto-ms folder, start fims backend service with `go run cmd/main.go`, running at port `5002`
3. in [fims-proto-ui](https://github.com/fims-proto/fims-proto-ui) folder, start vite in 'dev' mode with `npm run dev`, running at port `5001`
4. ory oathkeeper will perform as gateway and manage all traffic via port `4455`
4. open [`http://127.0.0.1:4455/ui/`](http://127.0.0.1:4455/ui/)
5. ready to go :)

## Test user
1. in [fims-dev-environment](https://github.com/fims-proto/fims-dev-environment) folder, run command `./scripts/user_creation.sh`, hit enter to confirm kratos admin API, then input user email address. E.g.:
``` shell
➜  fims-dev-environment git:(master) ./scripts/user_creation.sh
Kratos admin API [http://127.0.0.1:4434]:
User email address: tester@fims.com

User tester@fims.com created:
==> User id: 4e84575b-b5ba-493f-b5b6-dd2db0ef6fb0
```
2. note down created user id, e.g.: `4e84575b-b5ba-493f-b5b6-dd2db0ef6fb0`
3. run command `./scripts/user_invitation.sh`, input user id, hit enter. E.g.:
``` shell
➜  fims-dev-environment git:(master) ./scripts/user_invitation.sh
Kratos admin API [http://127.0.0.1:4434]:
User ID: 4e84575b-b5ba-493f-b5b6-dd2db0ef6fb0

Recovery link created:
==> Expires in: 30 mins
==> Follow link: http://127.0.0.1:4455/kratos/public/self-service/recovery?flow=d9ae1b0b-842e-4b7d-8571-57f63dcb7d9e&token=Nw0sIPrNzzOb7Vs1hw3e33qso5mhvHpY
```
4. click returned link, you will be directed to profile update page.
5. goto 更新密码 tab, enter your password and submit
6. ready to go :)

Note: if you forgot to save your password and closed the page, you need to re-invite user, like step 3.

## Swagger UI
After starting fims-proto-ms, visit [`http://127.0.0.1:4455/fims/swagger/index.html`](http://127.0.0.1:4455/fims/swagger/index.html).  

## Postman test
Import [postman collection](https://github.com/fims-proto/fims-proto-ms/tree/master/pm_collection) into Postman app
