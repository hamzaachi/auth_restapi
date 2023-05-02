package models

import (
	"fmt"
	"net/http"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserList struct {
	Users []User `json:"users"`
}

func (u *User) Bind(r *http.Request) error {
	if u.Username == "" || u.Password == "" {
		return fmt.Errorf("Username and Password are required")
	}

	return nil
}

func (*User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (*UserList) Render() error {
	return nil
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

func (*Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
