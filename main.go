package main

import (
    "database/sql"
    "log"
    "math/rand"
    "net/http"
    "time"

    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
)

// MyDB facilitates the addition of methods on top of a sql.DB.
type MyDB struct {
    *sql.DB
}

func main() {

    /* Seeding our random integer generator */
    rand.Seed( time.Now().UTC().UnixNano())

    /* db */
    sqldb, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/bettor")
    if err != nil {
        log.Fatal(err)
    }
    defer sqldb.Close()

    if err = sqldb.Ping(); err != nil {
        log.Fatal(err)
    }

    /* context */
    db := MyDB{ sqldb }

    /* router */
    r := mux.NewRouter()

    /* contacts */
    r.Methods("PUT","POST").Path("/contacts").HandlerFunc(db.ContactsHandler)

    /* verify */
    r.Methods("PUT","POST").Path("/verify").HandlerFunc(db.VerificationHandler)

    /* users */
    users := r.PathPrefix("/users").Subrouter()

    users.Methods("GET").Path("/{id}").HandlerFunc(db.UserShowHandler)
    users.Methods("PUT", "POST").Path("/{id}").HandlerFunc(db.UserUpdateHandler)
    users.Methods("DELETE").Path("/{id}").HandlerFunc(db.UserDeleteHandler)
    users.Methods("GET").Path("/{id}/bets").HandlerFunc(db.UserBetsHandler)
    users.Methods("GET").Path("/{id}/witnessing").HandlerFunc(db.UserWitnessingHandler)

    users.Methods("GET").HandlerFunc(db.UsersShowHandler)
    users.Methods("PUT", "POST").HandlerFunc(db.UsersCreateHandler)

    /* bets */
    bets := r.PathPrefix("/bets").Subrouter()

    bets.Methods("GET").Path("/{id}").HandlerFunc(db.BetShowHandler)
    bets.Methods("DELETE").Path("/{id}").HandlerFunc(db.BetDeleteHandler)
    bets.Methods("PUT", "POST").Path("/{id}/status").HandlerFunc(db.BetStatusHandler)

    bets.Methods("GET").HandlerFunc(db.BetsShowHandler)
    bets.Methods("PUT", "POST").HandlerFunc(db.BetsCreateHandler)

    /* serve */
    log.Println("Starting server on :8080")
    http.ListenAndServe(":8080", r)
}   