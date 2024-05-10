# SQL Boiler, Many-to-Many Relation

The problem with this demo is, that SQL Boiler does not generate the models for the SQL table `domain_admins`.
To reproduce it.

## Install SQL Boiler:
```shell
go install github.com/volatiletech/sqlboiler/v4@latest
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest
```
Create database structure with (make sure you fill you database settings in `.env` and `sqlboiler.toml`):

## Install SQL migration tool
For SQL migration I use golang-migrate, install it with:
```shell
scoop install migrate
```
## Apply migrations to database
Create a SQL migration file with:
```shell
make dbu
```

## Generate the go files with
```shell
go generate
```

## Test the models
```shell
go test -v ./internal/db/models
```