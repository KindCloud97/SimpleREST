package main

import (
	"encoding/json"
	"fmt"
	pb "github.com/KindCloud97/SimpleREST/proto"
	grpcserver "github.com/KindCloud97/SimpleREST/server"
	model "github.com/KindCloud97/SimpleREST/user"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

var Session *gocql.Session

//Additional requests for cassandra db
func getAllUsersCas() ([]model.User, error) {
	var user model.User
	users := make([]model.User, 0)
	rows := Session.Query(`SELECT * FROM simplerest.users`).Iter()
	for rows.Scan(&user.ID, &user.Email, &user.Firstname, &user.Lastname) {
		users = append(users, user)
	}
	return users, nil
}

func getUserCas(userId gocql.UUID) (model.User, error) {
	var user model.User
	err := Session.Query(`SELECT * FROM simplerest.users WHERE id = ?`,
		userId).Scan(&user.ID, &user.Email, &user.Firstname, &user.Lastname)
	if err != nil {
		return user, fmt.Errorf("cassadra SELECT error:%w", err)
	}
	return user, nil
}

func createUserCas(u *model.User) error {
	return Session.Query("INSERT INTO simplerest.users(id, email, firstname, lastname) VALUES(?, ?, ?, ?)",
		u.ID, u.Email, u.Firstname, u.Lastname).Exec()
}

func updateUserCas(u *model.User) error {

	return Session.Query("UPDATE simplerest.users SET email = ?, firstname = ?, lastname = ? WHERE id = ?",
		u.Email, u.Firstname, u.Lastname, u.ID).Exec()
}

func deleteUserCas(userId gocql.UUID) error {
	return Session.Query("DELETE FROM simplerest.users WHERE id = ?",
		userId).Exec()
}

//Common requests for db

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := getAllUsersCas()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Fatal(err)
	}
	if err = r.Body.Close(); err != nil {
		log.Fatal(err)
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	parseUUID, err := gocql.ParseUUID(params["id"])
	if err != nil {
		log.Fatal(err)
		return
	}

	user, err := getUserCas(parseUUID)
	if err != nil {
		http.Error(w, "No such user", http.StatusBadRequest)
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user model.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	user.ID, err = gocql.RandomUUID()
	if err != nil {
		log.Fatal(err)
	}

	err = createUserCas(&user)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	if err = r.Body.Close(); err != nil {
		log.Fatal(err)
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user model.User
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(r.Body)
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	err = updateUserCas(&user)
	if err != nil {
		http.Error(w, "No such user", http.StatusBadRequest)
		return
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	parseUUID, err := gocql.ParseUUID(params["id"])
	if err != nil {
		log.Fatal(err)
	}
	err = deleteUserCas(parseUUID)
	if err != nil {
		http.Error(w, "No such user", http.StatusBadRequest)
		return
	}
}

//CassandraDB

func NewCassandraDB() {
	var err error
	cluster := gocql.NewCluster("db")
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10

	//replace username and password fields with their real settings.
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "cassandra",
		Password: "cassandra",
	}
	for {
		Session, err = cluster.CreateSession()
		if err == nil {
			break
		}
		log.Println("\tsleeping fo 1 second...")
		time.Sleep(time.Second)
	}
	fmt.Println("\t > cassandra db initialized")

	// create keyspace
	err = Session.Query(
		"CREATE KEYSPACE IF NOT EXISTS" +
			" simplerest" +
			" WITH replication = {" +
			"'class': " +
			"'SimpleStrategy', " +
			"'replication_factor': 1" +
			"};",
	).Exec()
	if err != nil {
		log.Fatal("Cannot create a keyspace!", err)
		return
	}

	// create table
	err = Session.Query(
		"CREATE TABLE IF NOT EXISTS" +
			" simplerest.users" +
			" (" +
			"id uuid, " +
			"email text, " +
			"firstname text, " +
			"lastname text, " +
			"PRIMARY KEY (id)" +
			");",
	).Exec()
	if err != nil {
		log.Fatal("Cannot create a table!", err)
		return
	}
}

func main() {
	//REST API
	NewCassandraDB()
	r := mux.NewRouter()
	url := "/api/v1/users/"
	r.HandleFunc(url, GetAllUsers).Methods("GET")
	r.HandleFunc(url+"{id}", GetUser).Methods("GET")
	r.HandleFunc(url, CreateUser).Methods("POST")
	r.HandleFunc(url, UpdateUser).Methods("PUT")
	r.HandleFunc(url+"{id}", DeleteUser).Methods("DELETE")
	go func() {
		log.Fatal(http.ListenAndServe(":8000", r))
	}()

	//gRPC server
	s := grpc.NewServer()
	pb.RegisterCRUDServer(s, &grpcserver.GRPCServer{})

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error")
	}

	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
