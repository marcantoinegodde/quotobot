package server

import "encoding/gob"

type User struct {
	ID        string
	Email     string
	Name      string
	FirstName string
	LastName  string
	Avatar    string
}

func init() {
	gob.Register(User{})
}
