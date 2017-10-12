package main

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
	"fmt"
	"strings"
)

func (a *App) Login(w http.ResponseWriter, r *http.Request) {

	var user UserCredentials

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, err)
		return
	}

	userFull := User{
		UserName : user.Username,
		Password : user.Password,
	}
	userFull.login(a.DB)

	if userFull.Id == 0{
		w.WriteHeader(http.StatusForbidden)
		return
	}

	userFull.Token, err = GenerateRandomString(32)

	stringUser, err := json.Marshal(userFull)

	err = a.Redis.Set(userFull.Token, stringUser,0).Err()
	if err != nil {
		fmt.Print(err)
	}

	respondWithJSON(w, http.StatusOK, userFull)
}

func (a *App) ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")

	if len(splitToken) < 2 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	_, err := a.Redis.Get(splitToken[1]).Result()
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	next(w,r)

}

//register user
func (a *App) Registration(w http.ResponseWriter, r *http.Request)  {

}

//
func (a *App) SendToDialog(w http.ResponseWriter, r *http.Request) {

}

//get dialogs by user_id
func (a *App) GetDialogs(w http.ResponseWriter, r *http.Request)  {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	dialogs, err := getDialogs(a.DB, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, dialogs)
}


func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
