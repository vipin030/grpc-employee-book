# gRPC REST CRUD service in Go

Simple gRPC go service help to add employee info into PostgreSQL database :boom:

Technologies used :point_down:
1. Go 1.5
2. PostgreSQL 12.4 (pgx package)

## Installation

```
git clone https://github.com/vipin030/grpc-employee-book.git
cd grpc-employee-book
go mod download
```

## Running the app

First terminal:
```
go run cmd/server/main.go
```
Second terminal(Run the app as REST client)
```
go run cmd/rest-client/main.go
```
Third Terminal (Run the app as gRPC client)
```
go run cmd/grpc-client/main.go
```
