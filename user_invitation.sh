#!/bin/bash

if ! command -v jq &> /dev/null
then
    echo "`jq` not installed, exit..."
    exit
fi

read -p "Kratos admin API [http://127.0.0.1:4434]: " kratos_admin_api
kratos_admin_api=${kratos_admin_api:-http://127.0.0.1:4434}

read -p "User ID: " user_id
if [ -z "$user_id" ]
then
    echo "User ID is mandatory"
    exit
fi

echo

response=$(curl --request POST -sL \
    -w "%{http_code}" \
    --header "Content-Type: application/json" \
    --data '{
        "expires_in": "30m",
        "identity_id": "'$user_id'"
    }' \
    $kratos_admin_api/recovery/link
)

http_code=$(tail -n1 <<< "$response")
content=$(sed '$ d' <<< "$response")

if [ $http_code != "200" ]
then
    echo "[!] Create link request failed: "
    echo "$content"
    exit
fi

recovery_link=$(jq -r ".recovery_link" <<< "$content")
echo "Recovery link created:"
echo "==> Expires in: 30 mins"
echo "==> Follow link: $recovery_link"