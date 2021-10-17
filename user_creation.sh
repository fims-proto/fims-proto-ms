#!/bin/bash

if ! command -v jq &> /dev/null
then
    echo "`jq` not installed, exit..."
    exit
fi

read -p "Kratos admin API [http://127.0.0.1:4434]: " kratos_admin_api
kratos_admin_api=${kratos_admin_api:-http://127.0.0.1:4434}


read -p "User email address: " email_addr
if [ -z "$email_addr" ]
then
    echo "Email address is mandatory"
    exit
fi

echo

response=$(curl --request POST -sL \
    -w "%{http_code}" \
    --header "Content-Type: application/json" \
    --data '{
        "schema_id": "default",
        "traits": {
            "email": "'$email_addr'"
        }
    }' \
    $kratos_admin_api/identities
)

http_code=$(tail -n1 <<< "$response")
content=$(sed '$ d' <<< "$response")

if [ $http_code != "201" ]
then
    echo "[!] Create user request failed: "
    echo "$content"
    exit
fi

user_id=$(jq -r ".id" <<< "$content")
echo "User $email_addr created:"
echo "==> User id: $user_id"