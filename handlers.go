package main

import (
    "database/sql"
    "encoding/json"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

/* Handlers */
// ContactsHandler
// UsersShowHandler
// UsersCreateHandler

// UserShowHandler
// UserDeleteHandler
// UserBetsHandler
// UserWitnessingHandler

// BetsShowHandler
// BetsCreateHandler

// BetShowHandler
// BetDeleteHandler
// BetWinnerHandler

/* Basic StatusCode Responses */

// WriteError writes a JSON-formatted error response to a ResponseWriter.
func (w *http.ResponseWriter) WriteError(code int, errMsg string) {
    js, _ := json.Marshall(*GenerateError(code, errMsg))
    w.Write(js)
}

// WriteSuccess writes a JSON-formatted success response to a ResponseWriter.
func (w *http.ResponseWriter) WriteSuccess() {
    js, _ := json.Marshall(JSONResponse{ meta: Meta{ Code: 200 }})
    w.Write(js)
}
