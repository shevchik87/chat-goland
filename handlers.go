package main

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
)


func (a *App) Login(w http.ResponseWriter, r *http.Request)  {

}
func (a *App) Registration(w http.ResponseWriter, r *http.Request)  {

}
func (a *App) UserExist(w http.ResponseWriter, r *http.Request)  {

}

func (a *App) SendToDialog(w http.ResponseWriter, r *http.Request) {

}

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
