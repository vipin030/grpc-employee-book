package v1

import (
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jackc/pgx/v4"
	pb "github.com/vipin030/grpc-employee-book/proto"
	"golang.org/x/net/context"
)

type server struct{ Conn *pgx.Conn }

const (
	table = "employees"
)

func NewEmployeeServiceServer(conn *pgx.Conn) *server {
	return &server{Conn: conn}
}

func (s *server) CreateEmployee(ctx context.Context, em *pb.Employee) (*pb.ID, error) {
	if em.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "Name is empty, please try again")
	}
	insertQueryMeta := fmt.Sprintf(`insert into %v (name, city, salary) values ($1, $2, $3)`, table)
	_, err := s.Conn.Exec(context.Background(), insertQueryMeta, em.Name, em.City, em.Salary)

	return &pb.ID{Id: em.Id}, err
}

func (s *server) DeleteEmployee(ctx context.Context, em *pb.ID) (*pb.ID, error) {
	if em.Id < 1 {
		return nil, status.Error(codes.InvalidArgument, "ID is empty, please try again")
	}
	deleteQueryMeta := `delete from employees where id = $1`
	_, err := s.Conn.Exec(ctx, deleteQueryMeta, em.Id)

	return &pb.ID{Id: em.Id}, err
}

func (s *server) ReadEmployee(ctx context.Context, em *pb.ID) (*pb.Employee, error) {
	if em.Id < 1 {
		return nil, status.Error(codes.InvalidArgument, "ID is empty, please try again")
	}
	readQueryMeta := fmt.Sprintf(`select id,name,city,salary from %v where id = $1`, table)
	var a = &pb.Employee{}
	err := s.Conn.QueryRow(context.Background(), readQueryMeta, em.Id).Scan(&a.Id, &a.Name, &a.City, &a.Salary)
	return a, err
}

func (s *server) UpdateEmployee(ctx context.Context, em *pb.Employee) (*pb.ID, error) {
	if em.Id < 1 {
		return nil, status.Error(codes.InvalidArgument, "ID is empty, please try again")
	}
	readQueryMeta := fmt.Sprintf(`select id from %v where id = $1`, table)
	updateQueryMeta := fmt.Sprintf(`update %v set name = $1, city = $2 where id = $3`, table)
	var id int32
	err := s.Conn.QueryRow(context.Background(), readQueryMeta, em.Id).Scan(&id)
	if err != nil {
		log.Fatalf("Cannot find the Employee due to %v", err)
		return nil, status.Error(codes.InvalidArgument, "Employee does not found")
	}
	_, err = s.Conn.Exec(context.Background(), updateQueryMeta, em.Name, em.City, em.Id)
	return &pb.ID{Id: em.Id}, err
}
