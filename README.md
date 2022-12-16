# Night Owl

> Night owl is a boilerplate project which have goals to reduce overhead indirection and make it easier to navigate through the project by defining clearer separation of concern

## How to run

You can simply run the project using following command

```
go run main.go server-start
```

However, it's recommended to run this project using `docker-compose` with following command

```
docker-compose -f docker-compose.yml up --build
```

The app would available in `localhost:8000`, you can also set custom port by providing `APP_PORT` in `.env` file

## Migrations

We use golang-migrate, please install `brew install golang-migrate` to install with brew or you can see another instalation option found in their repostiory.

```
# Create migration file
migrate create -ext sql -dir migrations/sql -seq create_team_table

# Run migration
export POSTGRESQL_URL='postgres://root:pass@localhost:5432/ouroboros_db?sslmode=disable'
migrate -database ${POSTGRESQL_URL} -path migrations/sql up

migrate -database ${POSTGRESQL_URL} -path migrations/sql down
```

## Docs

We use swaggo to documented our endpoint, use these following command in root folder to generate specs

```
swag init -g internal/entry-point/rest/rest.go --parseDependency
```

Please refer to official repository for [annotation](https://github.com/swaggo/swag#declarative-comments-format)
