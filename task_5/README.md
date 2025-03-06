# TASK 5

---

Сервер отправляет сообщение при регистрации, обновлении, удалении пользователя

server --> notify-topic --> emailServiceMock

emailServiceMock --> notify-confirm --> server

emailServiceMock отправляет потверждение с паузой 10-30 сек

---

ENV
- APP_SALT
- APP_REST_PORT (default=8080)
- APP_RPC_PORT (default=50051)
- KAFKA_HOST (default=localhost)
- KAFKA_PORT (default=9092)

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
