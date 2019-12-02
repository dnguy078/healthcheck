# healthcheck
HTTP Server that executes periodic health checks for HTTP websites

## To Run
```
go run cmd/main.go
```

## To Test:
```
go test ./...
```

## Command Line Flags:
```
Defaults can be found in /cmd/main.go

--bind=127.0.0.1:8080
    Address to run http server on

--sslcert=certificate.crt --sslkey=certificate.key
    Pass in SSL Certs

--checkfrequency=30s
    Frequency healthchecks are performed

ie)
go run cmd/main.go --checkfrequency=1s
```


## GoDocs
```
[GoDocs](https://godoc.org/github.com/dnguy078/healthcheck)
```

## API:
### List Health Checks
Returns a list of health checks sorted by endpoint with paging support of 10 items per page. (pagination begins at 0, and sorted alphabetically by endpoint)
```json
Request:
curl http://127.0.0.1:8080/api/health/checks?page=0

Response:
{
    "items": [
        {
            "id": "C6C5B3DC-6685-7698-3CD5-C3AB7C10B3AC",
            "status": "200 OK",
            "code": 200,
            "endpoint": "https://www.blizzard.com/en-us/",
            "checked": 1574906832,
            "duration": "446.04934ms"
        },
    ],
    "page": 0,
    "total": 1,
    "size": 10
}
```

### Get Health Check
Return a single health check
```json
Request:
curl localhost:8080/api/health/checks/C6C5B3DC-6685-7698-3CD5-C3AB7C10B3AC

Response:
{
    "id": "C6C5B3DC-6685-7698-3CD5-C3AB7C10B3AC",
    "status": "200 OK",
    "code": 200,
    "endpoint": "https://www.blizzard.com/en-us/",
    "checked": 1574906993,
    "duration": "622.260455ms"
}
```

### Create Health Check

```json
Request:
curl -X POST http://localhost:8080/api/health/checks \
-d '{
    "endpoint":  "https://www.blizzard.com/en-us/"
}'

Response:
{
    "id": "95D87755-E3B9-66BE-549D-CB856EE71FCF",
    "endpoint": "https://www.blizzard.com/en-us/"
}
```

### Execute a Health Check
This will execute a health check, with a timeout provided in the query string

```json
curl -X POST "http://127.0.0.1:8080/api/health/checks/95D87755-E3B9-66BE-549D-CB856EE71FCF/try?timeout=10s"


Response:
{
    "id": "C6C5B3DC-6685-7698-3CD5-C3AB7C10B3AC",
    "status": "200 OK",
    "code": 200,
    "endpoint": "https://www.blizzard.com/en-us/",
    "checked": 1574907428,
    "duration": "442.441307ms"
}

// timed out

{
    "id": "C6C5B3DC-6685-7698-3CD5-C3AB7C10B3AC",
    "status": "Error",
    "code": 0,
    "endpoint": "https://www.blizzard.com/en-us/",
    "checked": 0,
    "duration": "54.436µs",
    "error": "Get https://www.blizzard.com/en-us/: context deadline exceeded"
}

```

### Delete a Health Check
Deletes a Healthcheck
```json
curl -X DELETE http://127.0.0.1:8080/api/health/checks/94a1d1e8-6e44-409e-9cb4-7bfcac2de1ae
```


