## 1. Run Postgresql migrations:
```shell
cd database
./migrate -database 'postgres://username:password@localhost:5432/database?sslmode=disable' -path migrations/psql up
```

## 2. Docker run nats jet-stream:
```shell
docker run -p 4222:4222 -ti nats:latest -js
```

## 3. Run redis:
```shell
redis-server
```

## 4. Run consumer:
```shell
cd consumer
go run cmd/app/main.go
```

## 5. Run api:
```shell
cd api
go run cmd/app/main.go
```

## 6. Run producer:
```shell
cd producer
go run .
```

## Bombardier stress test:
```shell
go install github.com/codesenberg/bombardier@latest
bombardier -c 100 -n 10000 http://localhost:8080/api/v1/orders/{exists_order_uid}
```

## Docker jaeger start
```shell
docker run -d --name jaeger \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 16686:16686 \
  -p 14286:14268/tcp \
  -p 4317:4317 \
  -p 4318:4318 \
  jaegertracing/all-in-one:latest
```