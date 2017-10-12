package main

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	_ "github.com/lib/pq"
	"fmt"
	"log"
	"net/http"
	"github.com/codegangsta/negroni"
	"github.com/go-redis/redis"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
	Redis  *redis.Client
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.Redis = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080"},
		AllowCredentials: true,
	})
	handler := c.Handler(a.Router)
	log.Fatal(http.ListenAndServe(":9000", handler))
}

func (a *App)initializeRoutes()  {
	a.Router.StrictSlash(true)
	a.Router.HandleFunc("/login", a.Login).Methods("POST")
	a.Router.HandleFunc("/register", a.Registration).Methods("POST")

	// Protected Endpoints
	a.Router.Handle("/dialogs/{id:[0-9]+}", negroni.New(negroni.HandlerFunc(a.ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(a.GetDialogs)),
	))
}