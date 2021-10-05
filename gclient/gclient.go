package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	pb "github.com/KindCloud97/SimpleREST/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"os"
	"strings"
)

var c pb.CRUDClient

func main() {
	url := flag.String("url", "", "")
	flag.Parse()

	if *url == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	conn, err := grpc.Dial(*url, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()


	c = pb.NewCRUDClient(conn)

	//if len(os.Args) >= 3 {
	//	calls(flag.Args())
	//} else {
		fmt.Println("input procedure")
		calls(input(">>"))
	//}
}

func calls(proc string) {
	var resp string
	var err error

	switch proc {
	case "getall":
		resp, err = getallUsers()
	case "get":
		resp, err = getUser()
	case "create":
		resp, err = createUser()
	case "update":
		resp, err = updateUser()
	case "delete":
		resp, err = deleteUser()
	default:
		os.Exit(0)
	}

	if err != nil {
		fmt.Println("\nError:", err.Error())
	} else {
		fmt.Println("\n", resp)
	}

	calls(input("\n>>"))
}

func getallUsers() (string, error) {
	resp, err := c.GetAllUsers(context.Background(), &emptypb.Empty{})

	if err != nil {
		return "", err
	}

	js, err:= json.Marshal(resp.Users)
	if err != nil {
		log.Fatal(err)
	}

	return string(js), nil
}

func getUser() (string, error) {
	id := input("id:")

	resp, err := c.GetUser(context.Background(), &pb.GetUserRequest{Id: id})

	if err != nil {
		return "", err
	}

	return resp.String(), nil
}

func createUser() (string, error) {
	user := &pb.User{
		Email:     input("email:"),
		Firstname: input("firstname:"),
		Lastname:  input("lastname:"),
	}

	resp, err := c.CreateUser(context.Background(), &pb.CreateUserRequest{User: user})

	if err != nil {
		return "", err
	}

	return resp.String(), nil
}

func updateUser() (string, error) {
	user := &pb.User{
		Id:        input("id:"),
		Firstname: input("firstname:"),
		Lastname:  input("lastname:"),
		Email:     input("email:"),
	}

	resp, err := c.UpdateUser(context.Background(), &pb.UpdateUserRequest{User: user})

	if err != nil {
		return "", err
	}

	return resp.String(), nil
}

func deleteUser() (string, error) {
	id := input("id:")

	resp, err := c.DeleteUser(context.Background(), &pb.DeleteUserRequest{Id: id})

	if err != nil {
		return "", err
	}

	return resp.String(), nil
}

func input(arg string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(arg + " ")
	input, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	input = strings.TrimSuffix(input, "\n")
	return input
}