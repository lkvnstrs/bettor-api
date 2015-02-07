package main

import (
    "database/sql"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

func main() {
    
    /* db */
    db, err := sql.Open("postgres", "postgresql://localhost/bettor")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    /* router */
    r := mux.NewRouter()

    /* contacts */
    r.HandlerFunc("/contacts")

    /* users */
    users := r.PathPrefix("/users").Subrouter()

    users.Methods("GET").HandlerFunc(db.UsersShowHandler)
    users.Methods("PUT", "POST").HandlerFunc(db.UsersCreateHandler)

    users.Methods("GET").Path("/{id}").HandlerFunc(db.UserShowHandler)
    users.Methods("DELETE").Path("/{id}").HandlerFunc(db.UserDeleteHandler)
    users.Methods("GET").Path("/{id}/bets").HandlerFunc(db.UserBetsHandler)
    users.Methods("GET").Path("/{id}/witnessing").HandlerFunc(db.UserWitnessingHandler)

    /* bets */
    bets := r.PathPrefix("/bets").Subrouter()

    bets.Methods("GET").HandlerFunc(db.BetsShowHandler)
    bets.Methods("PUT", "POST").HandlerFunc(db.BetsCreateHandler)

    bets.Methods("GET").Path("/{id}").HandlerFunc(db.BetShowHandler)
    bets.Methods("DELETE").Path("/{id}").HandlerFunc(db.BetDeleteHandler)
    bets.Methods("PUT", "POST").Path("/{id}/winner").HandlerFunc(db.BetWinnerHandler)
}   