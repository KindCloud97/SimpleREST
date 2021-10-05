// Package classification SimpleREST
//
// Documentation of our awesome API.
//
//     Schemes: http
//     BasePath: /api/v1/
//     Version: 1.0.0
//     Host: 127.0.0.1:8000
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//
// swagger:meta
package docs

import (
	"github.com/KindCloud97/SimpleREST/user"
)

// swagger:route GET /api/v1/users/{id} users getUser
// returns user by specified id
// Produces:
//     - application/json
// responses:
//   	200: UserGetIdResponse
//      400: badRequest

// swagger:parameters getUser
type UserIdParam struct {
	// Specifies uuid for a user
	//
	// unique: true
	// in: path
	// example: 3ca4ce84-ed71-42aa-8d1a-c0e001d3b8b4
	Id string `json:"id"`
}

// Error struct with error explanation string
// swagger:response badRequest
type BadRequestResponseWrapper struct {
	// in:body
	Body string
}


// swagger:response UserGetIdResponse
type UserGetIdResponse struct {
	// Specifies uuid for a user
	//
	// in: body
	//
	Body user.User
}


// swagger:route GET /api/v1/users/ users getUsers
// returns users
// Produces:
//     - application/json
// responses:
//   	200: UsersGetResponse
//      400: badRequest

// swagger:response UsersGetResponse
type UsersGetResponse struct {
	//in: body
	Body []user.User
}

// swagger:route POST /api/v1/users users createUser
// Produces:
//     - application/json
// responses:
//      400: badRequest

// swagger:parameters createUser
type UserPostParam struct {
	// in: body
	Body user.User
}

// swagger:route PUT /api/v1/users users updateUser
// Produces:
//     - application/json
// responses:
//      400: badRequest

// swagger:parameters updateUser
type UserPutParam struct {
	//in: body
	Body user.User
}

// swagger:route DELETE /api/v1/users/{id} users deleteUser
// Produces:
//     - application/json
// responses:
//      400: badRequest

// swagger:parameters deleteUser
type UserDelParam struct {
	// Specifies uuid for a user
	//
	// unique: true
	// in: path
	// example: 3ca4ce84-ed71-42aa-8d1a-c0e001d3b8b4
	Id string `json:"id"`
}