smock - configurable http-mock server

```
~[0]$ curl -s localhost:8080/_cnf | jq .
{
  "/bad": {
    "Code": 400,
    "Data": "{\"result\": \"BAD\"}"
  },
  "/fail": {
    "Code": 500,
    "Data": "{\"result\": \"FAIL\"}"
  },
  "/success": {
    "Code": 200,
    "Data": "{\"result\": \"OK\"}"
  }
}

~[0]$ curl -s localhost:8080/fail
{"result": "FAIL"}

~[0]$ curl -X POST localhost:8080/_cnf -d '{
  "path": "/abc",
  "code": 200,
  "data": "{\"x\": \"y\"}"
}
'
OK

~[0]$ curl -s localhost:8080/abc | jq .
{
  "x": "y"
}

~[0]$ curl -s localhost:8080/_cnf | jq .
{
  "/abc": {
    "Code": 200,
    "Data": "{\"x\": \"y\"}"
  },
  "/bad": {
    "Code": 400,
    "Data": "{\"result\": \"BAD\"}"
  },
  "/fail": {
    "Code": 500,
    "Data": "{\"result\": \"FAIL\"}"
  },
  "/success": {
    "Code": 200,
    "Data": "{\"result\": \"OK\"}"
  }
}

~[0]$ curl -X DELETE localhost:8080/_cnf -d '{
  "path": "/abc"
}
'
OK

 ~[0]$ curl -s localhost:8080/abc
404 page not found
```
