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
    sqldb, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/bettor?parseTime=true")
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

    users.Methods("GET").Path("/{id:[0-9]+}").HandlerFunc(db.UserShowHandler)
    users.Methods("PUT", "POST").Path("/{id:[0-9]+}").HandlerFunc(db.UserUpdateHandler)
    users.Methods("DELETE").Path("/{id:[0-9]+}").HandlerFunc(db.UserDeleteHandler)
    users.Methods("GET").Path("/{id:[0-9]+}/bets").HandlerFunc(db.UserBetsHandler)
    users.Methods("GET").Path("/{id:[0-9]+}/witnessing").HandlerFunc(db.UserWitnessingHandler)

    users.Methods("GET").HandlerFunc(db.UsersShowHandler)
    users.Methods("PUT", "POST").HandlerFunc(db.UsersCreateHandler)

    /* bets */
    bets := r.PathPrefix("/bets").Subrouter()

    bets.Methods("PUT", "POST").Path("/hook").HandlerFunc(db.BetsHookHandler)
    bets.Methods("GET").Path("/{id:[0-9]+}").HandlerFunc(db.BetShowHandler)
    bets.Methods("DELETE").Path("/{id:[0-9]+}").HandlerFunc(db.BetDeleteHandler)
    bets.Methods("PUT", "POST").Path("/{id:[0-9]+}/status").HandlerFunc(db.BetStatusHandler)

    bets.Methods("GET").HandlerFunc(db.BetsShowHandler)
    bets.Methods("PUT", "POST").HandlerFunc(db.BetsCreateHandler)

    r.Methods("OPTIONS").HandlerFunc(func (rw http.ResponseWriter, r *http.Request) {
        rw.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
        // rw.Header().Set("Access-Control-Content-Type", "*")
        rw.Header().Set("Access-Control-Allow-Origin", "*")
        rw.WriteHeader(200)
    })

    /* serve */
    log.Println("Starting server on :8080")
    http.ListenAndServe(":8080", r)
}   