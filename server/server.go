package grpcserver

import (
	"context"
	"fmt"
	pb "github.com/KindCloud97/SimpleREST/proto"
	"github.com/gocql/gocql"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

var db = new(gocql.Session)

type GRPCServer struct {}

func (s *GRPCServer) GetAllUsers(ctx context.Context, in *emptypb.Empty) (*pb.GetAllUsersResponse, error) {
	fmt.Println("\n\tOK")
	var users []*pb.User
	m := map[string]interface{}{}

	iter := db.Query("SELECT * FROM simplerest.users").Iter()
	for iter.MapScan(m) {
		users = append(users, &pb.User{
			Id:        m["id"].(gocql.UUID).String(),
			Email:     m["email"].(string),
			Firstname: m["firstname"].(string),
			Lastname:  m["lastname"].(string),
		})
		m = map[string]interface{}{}
	}

	return &pb.GetAllUsersResponse{Users: users}, nil
}

func (s *GRPCServer) GetUser(ctx context.Context, rq *pb.GetUserRequest) (*pb.User, error) {
	var user pb.User
	fmt.Println("\n\tOK")
	err := db.Query("SELECT id, email, firstname, lastname FROM simplerest.users WHERE id=?", rq.GetId()).
		Scan(&user.Id, &user.Email, &user.Firstname, &user.Lastname)

	return &user, err
}

func (s *GRPCServer) CreateUser(ctx context.Context, rq *pb.CreateUserRequest) (user *pb.User, err error) {
	user = rq.GetUser()
	fmt.Println("\n\tOK")
	uuid, err  := gocql.RandomUUID()
	if err != nil {
		log.Fatalln(err)
	}
	user.Id = uuid.String()

	err = db.Query("INSERT INTO simplrest.users(id, email, firstname, lastname) VALUES(?, ?, ?, ?)",
		user.GetId(), user.GetEmail(), user.GetFirstname(), user.GetLastname()).Exec()
	fmt.Println("\n\tOK2")
	return user, err
}

func (s *GRPCServer) UpdateUser(ctx context.Context, rq *pb.UpdateUserRequest) (user *pb.User, err error) {
	user = rq.GetUser()
	fmt.Println("\n\tOK")
	err = db.Query("UPDATE simplerest.users SET  email = ?, firstname = ?, lastname = ? WHERE id = ?",
		user.GetEmail(), user.GetFirstname(), user.GetLastname(), user.GetId()).Exec()

	return user, err
}

func (s *GRPCServer) DeleteUser(ctx context.Context, rq *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	fmt.Println("\n\tOK")
	if err := db.Query("DELETE FROM simplerest.users WHERE id = ?", rq.GetId()).Exec(); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}