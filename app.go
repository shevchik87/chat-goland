package main

import (
	"database/sql"

	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"log"
	"net/http"
	"github.com/shevchik87/chat-goland/socket"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
	Redis  *redis.Client
	Hub    *socket.Hub
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
	a.Hub = socket.NewHub()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	go a.Hub.Run()
	handler := c.Handler(a.Router)

	log.Fatal(http.ListenAndServe(":9000", handler))
}

func (a *App) initializeRoutes() {
	a.Router.StrictSlash(true)
	a.Router.HandleFunc("/login", a.LoginHandler).Methods("POST")
	a.Router.HandleFunc("/register", a.RegistrationHandler).Methods("POST")

	a.Router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		socket.ServeWs(a.Hub, w, r)
	})

	// Protected Endpoints
	a.Router.Handle("/dialogs/{id:[0-9]+}", negroni.New(negroni.HandlerFunc(a.ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(a.GetDialogsHandler)),
	))
	a.Router.Handle("/user/me", negroni.New(negroni.HandlerFunc(a.ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(a.MeHandler)),
	))
	//join to room
	//parameters user_id, hash room
	a.Router.Handle("/room/join", negroni.New(negroni.HandlerFunc(a.ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(a.JoinToRoomHandler)),
	)).Methods("POST")

	//create room
	//parameters name_room, type_room
	a.Router.Handle("/room/create", negroni.New(negroni.HandlerFunc(a.ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(a.CreateRoomHandler)),
	)).Methods("POST")

	//send message to room
	//parameter owner_Id room_id text
	a.Router.Handle("/message/send-to-room", negroni.New(negroni.HandlerFunc(a.ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(a.SendMessageToRoomHandler)),
	)).Methods("POST")

}
