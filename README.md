# Mindbody Scheduler

#### Mindbody-scheduler is a scheduler to run every alotted time to sign up for classes on [mindbody.io](mindbody.io) for Crossfit ECF only (currently)

_Built with Go as the server, Svelte as the client, and Postgres as the database._

## Run with [docker-compose](https://docs.docker.com/compose/)

```sh
> docker-compose build
> docker-compose up
```

## Run locally

_Requirements_

-   [go](https://golang.org/doc/go1.14) 1.14+
-   [npm](https://www.npmjs.com/package/npm) 6.9.0+
-   geckodriver (firefox) as the webdriver
-   java sdk 11+
-   xvfb

### server

```sh
> cd server
> GOOS=linux go build -o backend main.go
> ./backend &  # to run in background
```

### client

```sh
> cd web
> npm i
> npm run dev
```

### database

create local Postgres instance with username: `postgres` & password: `postgres`

```sh
> psql -h 127.0.0.1 -U postgres -d mb_scheduler_db -f scripts/init.sql
> psql -h 127.0.0.1 -U postgres -d mb_scheduler_db -c `Select * from schedule_rt`
```
