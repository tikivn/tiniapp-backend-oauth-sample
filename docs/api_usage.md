# API USAGE

## Get Access Token From Auth Code

```sh
curl -L -X POST 'http://localhost:8080/api/auth/token' -H 'Content-Type: application/json' --data-raw '{
  "auth_code": "<your-auth-code>"
}'
```

## Get Access Token From Refresh Token

```sh
curl -L -X POST 'http://localhost:8080/api/auth/token/refresh' -H 'Content-Type: application/json' --data-raw '{
  "refresh_token": "<your-refresh-token>"
}'
```

## Get Me From Access Token

```sh
curl -L -X POST 'http://localhost:8080/api/auth/me' -H 'Content-Type: application/json' --data-raw '{
  "access_token": "<your-access-token>"
}'
```
