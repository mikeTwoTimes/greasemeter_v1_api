## GreaseMeter API

A Monolithic REST API written in Go using the Gin web framework.

## Dependencies

- go1.24.9 or higher
- PostgreSQL 17 or higher

## Installation

```
git clone https://github.com/mikeTwoTimes/greasemeter_v1_api.git
cd greasemeter_v1_api
```

Before being able to run the server, you will need a .env file in the project 
root with the following sample values.

```
PORT=8080
JWT_SECRET=21429e89cc5fdea19d3d3d1073e3590a6d6c85c065c639eee29afe5d88d46b54
DB_CONN=postgresql://username:password@localhost:5432/grease_db
```

The server may run now, but it can't really do anything without a database 
connection. Run these commands from the project root to quickly set up a mock 
database.

```
createdb grease_db
psql grease_db < migrations/init.sql
psql grease_db < migrations/seed.sql
```

And be sure to update the username and password in the database connection 
environment variable!

## Running

```
go run ./cmd/api
```

The first time running the server, Go will automatically install all necessary 
packages. Subsequent runs will not require this step.

## Testing

Now that the server is up and running we can test it's functionality by running

```
go test ./tests -env local -v
```

from the project root. This will run a script that automates http requests to 
the server and determines whether or not they were successful by their response
codes. 

Please note that running more than one test in the span of a minute will
cause the server's rate limiter to kick in, resulting in a heap of 429 response
codes. Also, if your port environment variable is not 8080, you will need to
change line 18 in ./tests/api_test.go to match that value.

## Documentation

To see the complete documentation of the API, run the server and go to:
http://localhost:8080/swagger/index.html
