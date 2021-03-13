package main

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func init() {
  var err error
  Db, err = sql.Open("mysql", "root:root@tcp(db)/catechdojo")
  if err != nil {
    panic(err)
  }
}

// User
// Create a new user
func (user *User) create() (err error) {
  _, err = Db.Exec("insert into users (name, token) values (?, ?)", user.Name, user.Token)
  return
}
