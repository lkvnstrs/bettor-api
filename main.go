package main

import (
    "database/sql"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    /* Seeding our random integer generator */
    rand.Seed( time.Now().UTC().UnixNano())

    /* db */
    db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/bettor")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    if err = db.Ping(); err != nil {
        log.Fatal(err)
    }

    /* router */
    r := mux.NewRouter()

    /* contacts */
    r.HandlerFunc("/contacts", ContactsHandler)

    /* verify */
    r.HandlerFunc("/verify", VerificationHandler)

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
    bets.Methods("PUT", "POST").Path("/{id}/status").HandlerFunc(db.BetStatusHandler)

    /* serve */
    log.Println("Starting server on :8080")
    http.ListenAndServe(":8080", r)
}   