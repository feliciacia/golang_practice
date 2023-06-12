package main

import (
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID        string
	Username  string
	Email     string
	pswdHash  string
	CreatedAt string
	Active    string
	verHash   string
	timeout   string
}
