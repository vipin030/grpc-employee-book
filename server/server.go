package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jackc/pgx/v4"
	pb "github.com/vipin030/grpc-employee-book/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server struct{}

// Conn object
var Conn *pgx.Conn

func main() {
	var (
		username = "interface"
		password = "interface"
		host     = "localhost"
		schema   = "employee"
	)
	dsn := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", username, password, host, schema)
	fmt.Printf("Server gonna start ")
	lis, err := net.Listen("tcp", "127.0.0.1:50051")

	if err != nil {
		log.Fatalf("can't listen on port %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterEmployeeServiceServer(s, &server{})
	conn, err := pgx.Connect(context.Background(), dsn)
	Conn = conn

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	fmt.Println("Server has been started.. on ", lis.Addr().String())
	c := make(chan os.Signal)

	// os.Interrupt = CTRL+C
	signal.Notify(c, os.Interrupt)

	// Block main routine until a signal is received
	<-c

	// After receiving CTRL+C Properly stop the server
	fmt.Println("\nStopping the server...")
	s.Stop()
	lis.Close()
}

func (s *server) CreateEmployee(ctx context.Context, em *pb.Employee) (*pb.ID, error) {
	if em.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "Name is empty, please try again")
	}
	insertQueryMeta := `insert into employees (name, city, salary) values ($1, $2, $3)`
	_, err := Conn.Exec(context.Background(), insertQueryMeta, em.Name, em.City, em.Salary)

	return &pb.ID{Id: em.Id}, err
}

func (s *server) DeleteEmployee(ctx context.Context, em *pb.ID) (*pb.ID, error) {
	if em.Id < 1 {
		return nil, status.Error(codes.InvalidArgument, "ID is empty, please try again")
	}
	deleteQueryMeta := `delete from employees where id = $1`
	_, err := Conn.Exec(ctx, deleteQueryMeta, em.Id)

	return &pb.ID{Id: em.Id}, err
}

func (s *server) ReadEmployee(ctx context.Context, em *pb.ID) (*pb.Employee, error) {
	if em.Id < 1 {
		return nil, status.Error(codes.InvalidArgument, "ID is empty, please try again")
	}
	readQueryMeta := `select id,name,city,salary from employees where id = $1`
	var a = &pb.Employee{}
	err := Conn.QueryRow(context.Background(), readQueryMeta, em.Id).Scan(&a.Id, &a.Name, &a.City, &a.Salary)
	return a, err
}

func (s *server) UpdateEmployee(ctx context.Context, em *pb.Employee) (*pb.ID, error) {
	if em.Id < 1 {
		return nil, status.Error(codes.InvalidArgument, "ID is empty, please try again")
	}
	readQueryMeta := `select id from employees where id = $1`
	updateQueryMeta := `update employees set name = $1, city = $2 where id = $3`
	var id int32
	err := Conn.QueryRow(context.Background(), readQueryMeta, em.Id).Scan(&id)
	if err != nil {
		log.Fatalf("Cannot find the Employee due to %v", err)
		return nil, status.Error(codes.InvalidArgument, "Employee does not found")
	}
	_, err = Conn.Exec(context.Background(), updateQueryMeta, em.Name, em.City, em.Id)
	if err != nil {
		log.Fatalf("updation failed due to %v", err)
		return nil, status.Error(codes.InvalidArgument, "Updation failed")
	}
	return &pb.ID{Id: em.Id}, nil
}
