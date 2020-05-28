## Prerequisite
Only *Docker* is required to run the demo.
Postgresql server will be started in docker on port 5432, with db username: `postgres` and password: `fvpgsecret`
The local http server will be build and started in docker on port 3000, with basic auth username: `fv` and password: `fvpass`.

## Demo
```bash
# start the db first
./scripts/start-db.sh

# start the http server
./scripts/start.sh

# run the demo
# 1. create some data
curl -i -XPOST -d '{"name":"Apple", "valid": true, "count": 1}' http://fv:fvpass@localhost:3000/s/1
curl -i -XPOST -d '{"name":"Banana", "valid": false, "count": -12}' http://fv:fvpass@localhost:3000/s/1

# 2. do some query
curl -i -XGET http://fv:fvpass@localhost:3000/s/1/name/Apple

# 3. insert data
curl -i -XPOST -d '{"name":"Cookie", "valid": true, "count": 103}' http://fv:fvpass@localhost:3000/s/1

# 4. update data
# this is different from the requirement, need to do a query to get the rowId first
curl -i -XGET http://fv:fvpass@localhost:3000/s/1/name/Banana
# assume the id is 2, we can now do the update:
curl -i -XPOST -d '{"name":"Banana", "valid": true, "count": 12}' http://fv:fvpass@localhost:3000/s/1/r/2

```

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

POST /s/{specId}
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
Assume only query by single column and string is the only supported type.
The response is an array, since we have no way to specify a unique id.
Note: without column and value all rows will be returned.

GET /s/{specId}/{column}/{value}

example:
```
GET /s/1/name/Apple
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

PATCH /s/{specId}/r/{rowId}
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

POST /s/{specId}/r/{rowId}
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
