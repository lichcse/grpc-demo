1. Start server

- `make server` or `go run cmd/api/main.go`

2. Run client

- `make client` or `go run cmd/client/notification/main.go args`
  - args:
    - `SendMessage` : Unary RPCs
    - `ClientStreaming` : Client streaming RPCs
    - `ServerStreaming` : Server streaming RPCs
    - `BidirectionStreaming` : Bidirectional streaming RPCs
  - Ex: `go run cmd/client/notification/main.go SendMessage`

3. Health check with `grpcurl`

```
grpcurl -plaintext localhost:5001 demo.healthcheck.v2.HealthCheck/Check
```

4. Test send message with `grpcurl`

```
grpcurl -plaintext -d @ localhost:5001 demo.notication.v2.NoticationService/SendMessage <<EOM
{
  "message": "Hi"
}
EOM
```

5. Generate with prototool - File: `prototool.yaml`

- Edit path: `../../../pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.2.0/third_party/googleapis`
- `prototool generate`
