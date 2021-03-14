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

// Read a user
func retrieve(token string) (user User, err error) {
  user = User{}
  err = Db.QueryRow("select id, token, name from users where token = ?", token).Scan(&user.ID, &user.Token, &user.Name)
  return
}

// Update a user
func (user *User) update() (err error) {
  _, err = Db.Exec("update users set name = ?", user.Name)
  return
}

func (user *User) getAllCharacters() (userCharacterListResponses []UserCharacterListResponse, err error) {
  rows, err := Db.Query("select id, character_id from user_characters where user_id = ?", user.ID)
  for rows.Next() {
    userCharacterListResponse := UserCharacterListResponse{}
    err = rows.Scan(&userCharacterListResponse.ID, &userCharacterListResponse.CharacterID)
    if err != nil {
      return
    }
    var character Character
    character, err = retrive_character(userCharacterListResponse.CharacterID)
    if err != nil {
      return
    }
    userCharacterListResponse.Name = character.name
    userCharacterListResponses = append(userCharacterListResponses, userCharacterListResponse)
  }
  return
}

// UserCharacter
// Create a new user_character
func (uc *UserCharacter) create() (err error) {
  _, err = Db.Exec("insert into user_characters (user_id, character_id) values (?, ?)", uc.userID, uc.characterID)
  return
}

// Character
// Get all characters
func getAllCharacters() (characters []Character, err error) {
  rows, err := Db.Query("select id, name, weight from characters order by weight ASC")
  for rows.Next() {
    character := Character{}
    err = rows.Scan(&character.ID, &character.name, &character.weight)
    if err != nil {
      return
    }
    characters = append(characters, character)
  }
  return
}

// Get a Character
func retrive_character(cid int32) (character Character, err error) {
  err = Db.QueryRow("select name from characters where id = ?", cid).Scan(&character.name)
  return
}
