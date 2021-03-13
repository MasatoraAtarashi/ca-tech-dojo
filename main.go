package main

import (
  json2 "encoding/json"
  "github.com/dgrijalva/jwt-go"
  "github.com/gorilla/mux"
  "log"
  "math/rand"
  "net/http"
  "time"

  //"github.com/gorilla/sessions"
)

type User struct {
  Name string
  Token string
}

type UserCreateRequest struct {
  Name string `json:"name"`
}

type UserCreateResponse struct {
  Token string `json:"token"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
  len := r.ContentLength
  body := make([]byte, len)
  var userCreateRequest UserCreateRequest
  r.Body.Read(body)
  json2.Unmarshal(body, &userCreateRequest)
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

//func GetUser(w http.ResponseWriter, r *http.Request) {
//  //verifyKey, err := jwt.par
//  //TODO: メソッド切り出す
//  tokenString := r.Header.Get("x-token")
//  name := parseToken(tokenString)
//  if err != nil {
//    //TODO: 403を返す
//    return
//  }
//  //TODO: DB参照
//  //token, _ = jwt.Parse(myToken, func(token *jwt.Token) ([]byte, error) {
//  //  return myLookupKey(token.Header["kid"])
//  //})
//  fmt.Println(name)
//}

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

//func parseToken(tokenString string) (name string, err error) {
//  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//    return []byte("mySigningKey"), nil
//  })
//  if err != nil {
//    return
//  }
//  if token.Valid {
//    name = "a"
//    return
//  } else {
//    return
//  }
//}

func main() {
  r := mux.NewRouter()
  u := r.Path("/user").Subrouter()
  u.Methods("POST").HandlerFunc(CreateUser)
  //u.Methods("GET").HandlerFunc(GetUser)
  log.Fatal(http.ListenAndServe(":8080", r))
}
