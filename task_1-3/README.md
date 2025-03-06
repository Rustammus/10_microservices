# TASK 1-3

---
ENV
- APP_SALT
- APP_REST_PORT (default=8080)
- APP_RPC_PORT (default=50051)

## gRPC

Аутентификация только по gRPC

Proto спецификации ./proto/user.proto ./proto/auth.proto

## REST

- **POST /users**

```
{
    "name": "SomeName",
    "email": "email1",
    "password": "1234"
}
```

- **GET /users/{id}**

```

```

- **PUT /users/{id}**

```
{
    "name": "SomeName",
    "email": "email1",
    "password": "1234"
}
```

- **DELETE /users/{id}**

```

```
