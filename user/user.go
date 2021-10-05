package user

import (
	"github.com/gocql/gocql"
)

type User struct {
	ID      gocql.UUID `json:"id"`
	Email   string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
