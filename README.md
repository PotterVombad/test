чтобы запустить 
make up

чтобы получить пару
curl --location 'http://localhost:8080/tokens_pairs?uid=<uid>'

пример ответа 
{
    "access_token": "bla-bla",
    "refresh_token": "NGFmNjFiZmU3MDI2NDE0MmFhMWYxYmFjMzU3ODc2OGU="
}


curl --location --request PUT 'http://localhost:8080/refresh?token=refresh_token'