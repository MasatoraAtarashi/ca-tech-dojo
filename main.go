package main

import (
  json2 "encoding/json"
  "log"
  "math/rand"
  "net/http"
  "time"

  "github.com/dgrijalva/jwt-go"
  "github.com/gorilla/mux"
)

type User struct {
  Id int32
  Name string
  Token string
}

type UserCreateRequest struct {
  Name string `json:"name"`
}

type UserCreateResponse struct {
  Token string `json:"token"`
}

type UserGetResponse struct {
  Name string `json:"name"`
}

type UserUpdateRequest struct {
  Name string `json:"name"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
  len := r.ContentLength
  body := make([]byte, len)
  var userCreateRequest UserCreateRequest
  r.Body.Read(body)
  json2.Unmarshal(body, &userCreateRequest)
  //TODO: nameが指定されていなかったら400を返す
  name := userCreateRequest.Name

  token, err := createToken(name)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  user := User{Name: name, Token: token}
  err = user.create()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  userCreateResponse := UserCreateResponse{
    Token: token,
  }
  output, err := json2.Marshal(&userCreateResponse)
  w.Header().Set("Content-Type", "application/json")
  w.Write(output)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
  //TODO: x-tokenが指定されていなかったら400
  token := r.Header.Get("x-token")
  user, err := retrieve(token)
  if err != nil {
    http.Error(w, "不正なトークンです", http.StatusForbidden)
    return
  }
  userGetResponse := UserGetResponse{
    Name: user.Name,
  }
  output, err := json2.Marshal(&userGetResponse)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(output)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
  //TODO: x-tokenが指定されていなかったら400
  token := r.Header.Get("x-token")
  user, err := retrieve(token)
  if err != nil {
    http.Error(w, "不正なトークンです", http.StatusForbidden)
    return
  }

  len := r.ContentLength
  body := make([]byte, len)
  var userUpdateRequest UserUpdateRequest
  r.Body.Read(body)
  json2.Unmarshal(body, &userUpdateRequest)
  //TODO: nameが指定されていなかったら400を返す
  name := userUpdateRequest.Name

  user.Name = name
  err = user.update()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  return
}

func createToken(name string) (tokenString string, err error) {
  rand.Seed(time.Now().UnixNano())
  token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), jwt.MapClaims{
    "name": name,
    "createdAt": time.Now().UnixNano(),
    "rand": rand.Intn(9999999),
  })
  tokenString, err = token.SignedString([]byte("mySigningKey"))
  return
}

func main() {
  r := mux.NewRouter()
  u := r.Path("/user").Subrouter()
  u.Methods("POST").HandlerFunc(CreateUser)
  u.Methods("GET").HandlerFunc(GetUser)
  u.Methods("PUT").HandlerFunc(UpdateUser)
  log.Fatal(http.ListenAndServe(":8080", r))
}
