# Test API
## Run the following command to start the server:

```shell
go run cmd/server/main.go
```

## Open another tab in your terminal and run the following commands to add some records to your log:
```shell
curl -X POST localhost:8080 -d \
'{"record": {"value": "TGV0J3MgR28gIzEK"}}'
```
```shell
curl -X POST localhost:8080 -d \
'{"record": {"value": "TGV0J3MgR28gIzLK"}}'
```
```shell
curl -X POST localhost:8080 -d \
'{"record": {"value": "TGV0J3MgR28gIzMK"}}'
```

## Read records by following commands:
```shell
curl -X GET localhost:8080 -d '{"offset": 0}'
```
```shell
curl -X GET localhost:8080 -d '{"offset": 1}'
```
```shell
curl -X GET localhost:8080 -d '{"offset": 2}'
```

### Test not found record:
```shell
curl -X GET localhost:8080 -d '{"offset": 6969}'
```