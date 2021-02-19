package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/vipin030/grpc-employee-book/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "127.0.0.1:50051"
)

func main() {
	// Connect to the server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Can't connect to the server: %v", err)
	}

	defer conn.Close()

	c := pb.NewEmployeeServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Menu-based program allowing the user to choose from CRUD
	fmt.Println("\ngRPC CRUD Operations!")
	fmt.Print("Enter 1 => create an employee; 2 => List 3 => Modify and 4 => Remove: ")

	choice := bufio.NewReader(os.Stdin)
	text, _ := choice.ReadString('\n')
	text = strings.Trim(text, "\n")

	switch text {
	case "1":
		// Create Employee
		fmt.Print("\nEnter the name: ")
		nameBuf := bufio.NewReader(os.Stdin)
		name, _ := nameBuf.ReadString('\n')
		name = strings.Trim(name, "\n")

		fmt.Print("Enter the City: ")
		cityBuf := bufio.NewReader(os.Stdin)
		city, _ := cityBuf.ReadString('\n')
		city = strings.Trim(city, "\n")

		fmt.Print("Enter the Salary: ")
		salBuf := bufio.NewReader(os.Stdin)
		sal, _ := salBuf.ReadString('\n')
		sal = strings.Trim(sal, "\n")
		salInt, _ := strconv.Atoi(sal)

		_, err := c.CreateEmployee(ctx, &pb.Employee{Name: name, City: city,
			Salary: int32(salInt)})
		if err != nil {
			log.Fatalf("Failed to create new employee: %v", err)
		}
		fmt.Println("Saved Employee", "with the Name", name, "and salary", salInt)
	case "2":
		// List Employee
		fmt.Println("Enter the Employee ID")
		idBuf := bufio.NewReader(os.Stdin)
		id, _ := idBuf.ReadString('\n')
		id = strings.Trim(id, "\n")
		idInt, _ := strconv.Atoi(id)
		employees, err := c.ReadEmployee(ctx, &pb.ID{Id: int32(idInt)})
		if err != nil {
			log.Fatalf("Failed to read employee details due to %v", err)
		}
		fmt.Println("Employee details are below:\n", employees)
	case "3":
		// Update Employee
		fmt.Println("Enter the Employee ID")
		idBuf := bufio.NewReader(os.Stdin)
		id, _ := idBuf.ReadString('\n')
		id = strings.Trim(id, "\n")
		idInt, _ := strconv.Atoi(id)
		fmt.Println("Enter new Name: ")
		nameBuf := bufio.NewReader(os.Stdin)
		name, _ := nameBuf.ReadString('\n')
		name = strings.Trim(name, "\n")
		fmt.Println("Enter new City")
		cityBuf := bufio.NewReader(os.Stdin)
		city, _ := cityBuf.ReadString('\n')
		city = strings.Trim(city, "\n")
		_, err := c.UpdateEmployee(ctx, &pb.Employee{Id: int32(idInt), Name: name, City: city})
		if err != nil {
			log.Fatalf("Failed to update employee due to %v", err)
		}
		fmt.Println("Employee details has been updated")
	case "4":
		// Delete Employee
		fmt.Println("Enter the Employee ID")
		idBuf := bufio.NewReader(os.Stdin)
		id, _ := idBuf.ReadString('\n')
		id = strings.Trim(id, "\n")
		idInt, _ := strconv.Atoi(id)
		employee, err := c.DeleteEmployee(ctx, &pb.ID{Id: int32(idInt)})
		if err != nil {
			log.Fatalf("Failed to delete employee with Id %v due to %v", idInt, err)
		}
		fmt.Printf("Employee with ID %v has been deleted \n", employee.Id)

	default:
		fmt.Println("\nWrong option!")
	}
}
