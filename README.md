to run aplication

```make up```


to recive token pairs
```curl --location 'http://localhost:8080/tokens_pairs?uid=<uid>'```

response example
```
{
    "access_token": "access-token",
    "refresh_token": "NGFmNjFiZmU3MDI2NDE0MmFhMWYxYmFjMzU3ODc2OGU="
}```


curl --location --request PUT 'http://localhost:8080/refresh?token=refresh_token'