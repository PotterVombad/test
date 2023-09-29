to run aplication

```make up```


to recive token pairs

```curl --location 'http://localhost:8080/tokens_pairs?uid=<uid>'```

response example
```
{
    "access_token": "access-token",
    "refresh_token": "refresh-token"
}
```


to refresh tokens

```curl --location --request PUT 'http://localhost:8080/refresh?token=<refresh-token>'```


environment:

MONGO_PASSWORD\
MONGO_ADDR\
MONGO_USERNAME\
MONGO_DB_NAME\
MONGO_COL\
JWT_SECRET_KEY