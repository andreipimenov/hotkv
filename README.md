## Hot Key-Value In-Memory Storage

### Task

Implement hot key-value storage in memory without DB on build-in map type.
Requirements:
1. Managing storage with only one service.
2. Any read/write operations. Implement concurrent access model to the data.
3. Access to the data must be thread-safe.
4. Once key is got successfully it should be deleted from storage.
5. Unread data should be deleted from storage after 30 sec timeout.
6. Deleting should being proceed with channels not interating in order to search expired data.
7. Huge advantage is to use context.WithTimeout not time.After.

### Example of usage

Build service and run.
```
cd cmd/server
go build -o server
./server
```

Service provides http-api for managing storage.

Health check
```
curl -X GET 127.0.0.1:8080/api/ping
```
```
{
    "code": "OK",
    "message": "Pong"
}
```

Set key
```
curl -X POST -d '{"key":"hello", "value":"world"}' 127.0.0.1:8080/api/keys
```
```
{
    "code": "Created",
    "message": "Key hello is created successfully"
}
```

Get key
```
curl -X GET 127.0.0.1:8080/api/keys/hello
```
```
{
    "key": "hello",
    "value": "world"
}
```

Get key which was got previosly
```
{
    "code": "Not Found",
    "message": "key hello not found"
}
```

### Testing
Functional testing
```
go test -count=1 -cpu=1 -v ./storage
```

```
=== RUN   TestNew
--- PASS: TestNew (0.00s)
=== RUN   TestSetGet
--- PASS: TestSetGet (1.40s)
PASS
ok      github.com/andreipimenov/hotkv/storage  2.123s
```

Benchmarks
```
go test -count=1 -cpu=1 -bench=. -run=^Benchmark -benchmem ./storage
```
```
pkg: github.com/andreipimenov/hotkv/storage
BenchmarkSetKey           100000             13124 ns/op            1071 B/op          9 allocs/op
BenchmarkGetKey          3000000               517 ns/op              67 B/op          3 allocs/op
BenchmarkSetGetKey       1000000              4099 ns/op             258 B/op          7 allocs/op
PASS
ok      github.com/andreipimenov/hotkv/storage  11.650s
```