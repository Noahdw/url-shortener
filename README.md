
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
│  ├─ handler
│  └─ repository
│     ├─ db.go
│     ├─ models.go
│     └─ query.sql.go
└─ sqlc.yml

```
Url shortener with postgres db, sqlc for code generation, swagger api docs, docker compose with Watch for hot-reloading.
