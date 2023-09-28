start db linux:
```shell
sudo docker run --name=phatic_dialogue -e POSTGRES_PASSWORD=123456 -p   7766:5432 -d --rm postgres
```

start db windows:
```shell
docker run --name=phatic_dialogue -e POSTGRES_PASSWORD=123456 -p   7766:5432 -d --rm postgres
```

create db:
```shell
psql postgres://postgres:123456@localhost:7766
create database phatic_dialogue;
\q
```
run psql:
```shell
psql postgres://postgres:123456@localhost:7766
```
create db in psql mode:
```psql
create database phatic_dialogue;
```
drop db in psql mode:
```psql
drop database phatic_dialogue;
```
log out from psql mode:
```psql
\q
```

run seed:
```shell
go run cmd/main.go seed
```

run application
```shell
go run cmd/main.go run
```