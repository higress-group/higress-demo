编译：

```
go mod tidy
tinygo build -o main.wasm -scheduler=none -target=wasi -gc=custom -tags=custommalloc main.go
```

运行：

```
mv main.wasm local/
cd local
docker-compose up  # Ubuntu
docker compose up  # Mac
```

测试：
```
curl localhost:8080/ip
```

envoy log:
```
local-envoy-logs-1  | [2023-09-13 09:19:27.818983][35][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1188] wasm log httpcall-filter wasm-001: [http-call] get status: 200, request body: {
local-envoy-logs-1  |   "code": 200,
local-envoy-logs-1  |   "data": {
local-envoy-logs-1  |     "message": "ok"
local-envoy-logs-1  |   },
local-envoy-logs-1  |   "HttpStatusCode": 200,
local-envoy-logs-1  |   "successResponse": 200
local-envoy-logs-1  | }
```