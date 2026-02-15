# Lotto Journal

## Set up docker and database migration

- copy `.env.example` to `.env.local` and change values (if needed)
- Run

```shell
docker-compose --env-file .env.local up -d
```

- test connection on DB explorer program with connection string
- cd to /apps/api and run

```shell
migrate -path migrations -database "postgresql://${DB_USERNAME}:${DB_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=disable" up
```
