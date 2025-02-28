```
url-shortener
├─ api
├─ cmd
│  └─ api
│     ├─ docs
│     │  ├─ docs.go
│     │  ├─ swagger.json
│     │  └─ swagger.yaml
│     └─ main.go
├─ db
│  ├─ queries
│  │  └─ query.sql
│  └─ schema
│     └─ schema.sql
├─ docker-compose.yml
├─ Dockerfile
├─ go.mod
├─ go.sum
├─ internal
│  ├─ app
│  │  └─ app.go
│  ├─ handler
│  │  └─ url.go
│  ├─ repository
│  │  ├─ db.go
│  │  ├─ dbmock.go
│  │  ├─ models.go
│  │  └─ query.sql.go
│  └─ service
│     ├─ errors.go
│     ├─ url.go
│     └─ url_test.go
├─ README.md
└─ sqlc.yml

```

Url shortener with postgres db, sqlc for code generation, swagger api docs, docker compose with Watch for hot-reloading.

Architecture broken down into handler layer (for http), service layer (business logic), and repo (db) - enabling unit tests via interfaces.
