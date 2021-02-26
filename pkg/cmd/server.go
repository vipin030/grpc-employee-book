package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"

	"github.com/vipin030/grpc-employee-book/pkg/protocol/grpc"
	"github.com/vipin030/grpc-employee-book/pkg/protocol/rest"
	v1 "github.com/vipin030/grpc-employee-book/server"
)

// Config is configuration for Server
type Config struct {
	// gRPC server start parameters section
	// GRPCPort is TCP port to listen by gRPC server
	GRPCPort string

	// HTTP/REST gateway start parameters section
	// HTTPPort is TCP port to listen by HTTP/REST gateway
	HTTPPort string

	// DB Datastore parameters section
	// DBHost is host of database
	DBHost string
	// DBUser is username to connect to database
	DBUser string
	// DBPassword password to connect to database
	DBPassword string
	// DBSchema is schema of database
	DBSchema string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// get configuration
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "50051", "gRPC port to bind")
	flag.StringVar(&cfg.HTTPPort, "http-port", "8080", "HTTP port to bind")
	flag.StringVar(&cfg.DBHost, "db-host", "localhost", "Database host")
	flag.StringVar(&cfg.DBUser, "db-user", "interface", "Database user")
	flag.StringVar(&cfg.DBPassword, "db-password", "interface", "Database password")
	flag.StringVar(&cfg.DBSchema, "db-schema", "employee", "Database schema")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	if len(cfg.HTTPPort) == 0 {
		return fmt.Errorf("invalid TCP port for HTTP gateway: '%s'", cfg.HTTPPort)
	}

	dsn := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBSchema)
	conn, err := pgx.Connect(context.Background(), dsn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	v1API := v1.NewEmployeeServiceServer(conn)

	// run HTTP gateway
	go func() {
		_ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
	}()

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
