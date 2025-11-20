## GreaseMeter API

A Monolithic REST API written in Go using the Gin web framework.

## Dependencies

- go1.24.9 or higher

## Installation

```
git clone https://github.com/mikeTwoTimes/greasemeter_v1_api.git
cd greasemeter_v1_api
```

Before being able to run the server, you will need a .env file in the project 
root with the following values.

```
PORT=8080
JWT_SECRET=
DB_CONN=
SENDGRID_KEY=
```

The server will run just fine with those values. Obviously any operation 
requiring a server secret, database connection, or sendgrid key will fail if 
it's field is left blank.

## Running

```
go run ./cmd/api
```

The first time running the server, Go will automatically install all necessary 
packages. Subsequent runs will not require this step.

## Testing

To test our hosted API, you can run

```
go test -v -count=1 ./tests
```

from the project root. This will run a script that automatically tests our 
endpoints by comparing expected status codes to the ones received. Make note 
that running this script back to back within the span of a minute will cause 
the servers rate limiter to kick in, which will result in a heap of 429 status 
codes.

## Documentation

The hosted API's documentation can be found 
[here](https://api.greasemeter.live/swagger/index.html). 
