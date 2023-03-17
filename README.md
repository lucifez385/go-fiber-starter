# go-fiber-starter

## How to start project

```sh
docker compose up
```

## Healthz And ARGS Command

สามารถเพิ่ม args เพิ่มเติมเพื่อ run script บางอย่างได้ เช่น

```sh
go run server.go healthz
```

command name สามารถเพิ่มได้ทีี่ `server.go -> handleArgs()`

