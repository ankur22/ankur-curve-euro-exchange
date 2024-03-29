# Ankur's Exchange Rate Server Service for Curve

This will look at the exchange of currencies between EUR, USD and GBP:

1. Get today's exchange value from currency1 to currency2.
2. Get last week's exchange value from currency1 to currency2.
3. Perform a naive comparison between the two to determine if today is a good time to exchange currency1 into currency2.
4. Return the results -- single unit of curreny1 to curreny2 and the result of the naive comparison.

## API

### Request - `/v1/exchange?from={"EUR", "USD", "GBP"}&to={"EUR", "USD", "GBP"}`
Type: `GET`
<br />
<br />
Query parameter: `from` (required)
<br />
Valid values: {"EUR", "USD", "GBP"}
<br />
<br />
Query parameter: `to` (required)
<br />
Valid values: {"EUR", "USD", "GBP"}

#### Response
Status: `200`
<br />
Body: `{"dataDateTime":"2019-10-14T19:21:48.11587894+01:00","from":"EUR","shouldExchange":false,"singleUnit":1.1031,"to":"USD"}`
<br />
<br />
Status: `400`
<br />
Body: `{"reason":"query params are invalid. EUR, USD and GBP are valid."}`
<br />
<br />
Status: `500`
<br />
Body: `{"reason":""}`

#### cURL

Request:
<br />

```
curl http://localhost:8080/v1/exchange?from=EUR\&to=USD -v
```

<br />
Response:
<br />

```
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /v1/exchange?from=EUR&to=USD HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.58.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Mon, 14 Oct 2019 18:28:54 GMT
< Content-Length: 122
< 
{"dataDateTime":"2019-10-14T19:28:54.823492804+01:00","from":"EUR","shouldExchange":false,"singleUnit":1.1031,"to":"USD"}
* Connection #0 to host localhost left intact
```

## Build

The following will build and place a binary file in the `release/1.0.0/` directory:

```
./build.sh
```

## Run

Run with:

```
./release/1.0.0/exchange-1.0.0
```

## Tests

To run all tests:

```
go test ./...
```

## Bugs and Improvements

1. 500 response with no reason in json response body.
2. Requires logging to be implemented.
3. Requires [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper) integration.
4. More unit tests around failure cases, especially when a thread is getting data from network and other threads have to use stale data.
5. Integration tests.

## Test Environment

This was tested and working on:

 - Lenovo ThinkPad X1 Carbon - Core i7-4600U - 8 GB Ram
 - Ubuntu 18.04.3 LTS (binoic)
