# zcache

## Compile proto file

```
protoc --go_out=./zcachepb zcachepb/zcachepb.proto
```

## Run demo

```
# run on shell 1
go run main.go --addr="http://localhost:8001" --peers="http://localhost:8001,http://localhost:8002,http://localhost:8003" --api

# run on shell 2
go run main.go --addr="http://localhost:8002" --peers="http://localhost:8001,http://localhost:8002,http://localhost:8003"

# run on shell 3
go run main.go --addr="http://localhost:8003" --peers="http://localhost:8001,http://localhost:8002,http://localhost:8003"

# run on shell 4
curl "localhost:9999/api?key=zhangsan"
```
