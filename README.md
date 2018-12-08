# prest-dumpdata

## Command

```
go run main.go [action] [filename]
```
Action = dumpdata OR loaddata
Default action = dumpdata
Default filename = fixture.json

## Examples

```
go run main.go dumpdata test.json
```

```
go run main.go loaddata test.json
```

## Test database

- Up the container with postgresql
```
docker-compose up prest-db
```

- Copy sqls into container
```
sudo cp testdata/*.sql data/postgres/
```

- Drop database
```
docker-compose exec prest-db psql -U postgres -c "DROP DATABASE IF EXISTS prest;"
```

- Create database
```
docker-compose exec prest-db psql -U postgres -c "CREATE DATABASE prest;"
```

- Populate database
```
docker-compose exec prest-db psql -U postgres -d prest -1 -f /var/lib/postgresql/data/populate.sql
```

- Clears database
```
docker-compose exec prest-db psql -U postgres -d prest -1 -f /var/lib/postgresql/data/truncate.sql
```