package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Email:    "test@mail.ru",
		Password: "password",
	}
}
