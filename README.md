## Authentication
Basic Auth is used for simplicity, use username: `fv` and password: `fvpass`.

```bash
curl -i http://fv:fvpass@localhost:3000
HTTP/1.1 200 OK
Date: Wed, 27 May 2020 06:47:41 GMT
Content-Length: 7
Content-Type: text/plain; charset=utf-8
```

Without auth, http error 401 is returned.
```bash
curl -i http://localhost:3000
HTTP/1.1 401 Unauthorized
Www-Authenticate: Basic realm="fv"
Date: Wed, 27 May 2020 06:41:16 GMT
Content-Length: 0
```


## API
### Insert
*Warning* Not idempotent. If we want idempotent we can included a idempotent key in the header, or the call need to provide a uniqueID from caller.

POST /d/{id}
```json
{
    "name": "Cookie",
    "valid": true,
    "count": 103
}
```

response:
```json
{
    "id": "3",
    "name": "Cookie",
    "valid": true,
    "count": 103
}
```

### Query
Assume only query by single column.
The response is an array, since we have no way to specify a unique id.
Note: without column and value all rows will be returned.

GET /d/{id}/{column}/{value}

example:
```
GET /d/1/name/Apple
```

response:
```json
[
    {
        "id": "1",
        "name": "Apple",
        "valid": true,
        "count": 1 
    }
]
```

### Modify (not implemented)
This patches the row with changes.
This is not implemented as there would need additional api to unset a field.
Instead see below POST method, where a full row is updated.

PATCH /d/{id}/r/{row_id}
```json
{
    "valid": true,
}
```

response:
```json
{
    "id": "2",
    "name": "Banana",
    "valid": true,
    "count": 12
}
```

### Modify the row
*Warning* Not concurrency safe.
We could implement optimistic lock for concurrency, i.e. include a "version" field. But this is out of scope.

POST /d/{id}/r/{row_id}
```json
{
    "name": "Banana",
    "valid": true,
    "count": 12 
}
```

response:
```json
{
    "id": "2",
    "name": "Banana",
    "valid": true,
    "count": 12
}
