package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

func (a *App) LoginHandler(w http.ResponseWriter, r *http.Request) {

	var user UserCredentials

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, err)
		return
	}

	userFull := User{
		UserName: user.Username,
		Password: user.Password,
	}
	userFull.login(a.DB)

	if userFull.Id == 0 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	userFull.Token, err = GenerateRandomString(32)

	token := Token{
		Token: userFull.Token,
	}
	stringUser, err := json.Marshal(userFull)

	err = a.Redis.Set(userFull.Token, stringUser, 0).Err()
	if err != nil {
		fmt.Print(err)
	}

	respondWithJSON(w, http.StatusOK, token)
}

func (a *App) ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	token, err := extractToken(r)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	_, err = a.Redis.Get(token).Result()
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	next(w, r)

}

//register user
func (a *App) RegistrationHandler(w http.ResponseWriter, r *http.Request) {

}

//get user info
func (a *App) MeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := extractToken(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Token")
		return
	}

	user := User{}
	err = user.me(a, token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No data user")
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

//get dialogs by user_id
func (a *App) GetDialogsHandler(w http.ResponseWriter, r *http.Request) {

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

func (a *App) JoinToRoomHandler(w http.ResponseWriter, r *http.Request) {

	var room_user RoomUser
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&room_user); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := room_user.join(a.DB); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	respondWithJSON(w, http.StatusCreated, room_user)

}
func (a *App) CreateRoomHandler(w http.ResponseWriter, r *http.Request) {

	var room Room
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&room); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := room.create(a.DB); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	respondWithJSON(w, http.StatusCreated, room)

}
func (a *App) SendMessageToRoomHandler(w http.ResponseWriter, r *http.Request) {

	var message RoomMessage

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&message); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := message.Send(a.DB); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	respondWithJSON(w, http.StatusCreated, message)

}
func (a *App) SendMessageToDialogHandler(w http.ResponseWriter, r *http.Request) {

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

//function extract token from header
//return error or key
func extractToken(r *http.Request) (string, error) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")

	if len(splitToken) < 2 {
		return "", errors.New("No token")
	}

	return splitToken[1], nil
}
