package main

import (
  json2 "encoding/json"
  "fmt"
  "log"
  "math/rand"
  "net/http"
  "time"

  "github.com/dgrijalva/jwt-go"
  "github.com/gorilla/mux"
)

var error_map = map[string]int{
  "トークンを指定してください": http.StatusBadRequest,
  "不正なトークンです": http.StatusForbidden,
}

type MyError struct {
  msg string
}

func (e *MyError) Error() string {
  return fmt.Sprintf("%s", e.msg)
}

type User struct {
  ID int32
  Name string
  Token string
}

type Character struct {
  ID int32
  name            string
  weight int32
}

type UserCharacter struct {
  ID int32
  userID int32
  characterID     int32
  name            string
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

type GachaDrawResponse struct {
  Results []GachaResult `json:"results"`
}

type GachaDrawRequest struct {
  Times int32 `json:"times"`
}

type GachaResult struct {
  CharacterID int32 `json:"characterID"`
  Name string `json:"name"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
  len := r.ContentLength
  body := make([]byte, len)
  var userCreateRequest UserCreateRequest
  r.Body.Read(body)
  json2.Unmarshal(body, &userCreateRequest)
  name := userCreateRequest.Name
  if name == "" {
    http.Error(w, "名前を指定してください", http.StatusBadRequest)
    return
  }

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
  user, err := authenticate(r)
  if err != nil {
    http.Error(w, err.Error(), error_map[err.Error()])
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
  user, err := authenticate(r)
  if err != nil {
    http.Error(w, err.Error(), error_map[err.Error()])
    return
  }

  len := r.ContentLength
  body := make([]byte, len)
  var userUpdateRequest UserUpdateRequest
  r.Body.Read(body)
  json2.Unmarshal(body, &userUpdateRequest)
  name := userUpdateRequest.Name
  if name == "" {
    http.Error(w, "名前を指定してください", http.StatusBadRequest)
    return
  }

  user.Name = name
  err = user.update()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}

func DrawGacha(w http.ResponseWriter, r *http.Request) {
  user, err := authenticate(r)
  if err != nil {
    http.Error(w, err.Error(), error_map[err.Error()])
    return
  }

  len := r.ContentLength
  body := make([]byte, len)
  var gachaDrawRequest GachaDrawRequest
  r.Body.Read(body)
  json2.Unmarshal(body, &gachaDrawRequest)
  times := gachaDrawRequest.Times

  var results []GachaResult
  for ; times>0; times-=1 {
    gachaResult, err := weightPick()
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    userCharacter := UserCharacter{
     userID: user.ID,
     characterID: gachaResult.CharacterID,
    }
    userCharacter.create()

    results = append(results, gachaResult)
  }

  gachaDrawResponse := GachaDrawResponse{
    Results: results,
  }
  output, err := json2.Marshal(&gachaDrawResponse)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(output)
}

func authenticate(r *http.Request) (user User, err error) {
  token := r.Header.Get("x-token")
  if token == "" {
    err = &MyError{msg: "トークンを指定してください"}
    return
  }

  user, err = retrieve(token)
  if err != nil {
    err = &MyError{msg: "不正なトークンです"}
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

func weightPick() (gachaResult GachaResult, err error) {
  //characterを全部取得(重み昇順でソート)
  characters, err := getAllCharacters()
  if err != nil {
    return
  }

  //重みを合計する
  total_weight, err := sum_weight(characters)
  if err != nil {
    return
  }

  //乱数生成
  rand.Seed(time.Now().UnixNano())
  rnd := rand.Intn(int(total_weight))

  //生成された数字に基づいて返すcharacterを決める
  var picked Character
  for i:=0; i<len(characters); i+=1 {
    if rnd < int(characters[i].weight) {
      picked = characters[i]
      break
    }
    rnd -= int(characters[i].weight)
  }

  gachaResult = GachaResult{CharacterID: picked.ID, Name: picked.name}
  return
}

func sum_weight(characters []Character) (total_weight int32, err error) {
  for i:=0; i<len(characters); i+=1 {
    total_weight += characters[i].weight
  }
  return
}

func main() {
  r := mux.NewRouter()
  u := r.Path("/user").Subrouter()
  u.Methods("POST").HandlerFunc(CreateUser)
  u.Methods("GET").HandlerFunc(GetUser)
  u.Methods("PUT").HandlerFunc(UpdateUser)
  r.Path("/gacha/draw").Methods("POST").HandlerFunc(DrawGacha)
  log.Fatal(http.ListenAndServe(":8080", r))
}
