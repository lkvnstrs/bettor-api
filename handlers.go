package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "strconv"

    _ "github.com/go-sql-driver/mysql"
)

/* Handlers */

// ContactsHandler handles parsing contacts for existing users.
// Handles PUT and POST to /contacts.
func (db *sql.DB) ContactsHandler(rw http.ResponseWriter, r *http.Request) {

    var contactpairs []ContactPair
    var err error

    // parse the form
    phonenumbers := r.FormValue("phone_numbers")

    contactpairs, err = db.CheckPhoneNumbers(phonenumbers[0:])
    if err != nil {
        rw.WriteError(500, err.Error())
        return
    }

    // form as a JSON response
    m := Meta { Code: 200 }
    resp := JSONResponse { meta: m, data: contactpairs }

    // marshall and write
    js, err := json.Marshall(*resp)
    if err != nil {
        rw.WriteError(500, err.Error())
        return
    }

    rw.Write(js)
}

// VerificationHandler handles the verification of a user's phone number.
// Handles POST to /verify.
func (db *sql.DB) VerificationHandler(rw http.ResponseWriter, r *http.Request) {

}

// UsersShowHandler handles display of users.
// Handles GET to /users.
func (db *sql.DB) UsersShowHandler(rw http.ResponseWriter, r *http.Request) {

    var users []User 
    var err error

    // parse the form
    if err = r.ParseForm(); err != nil {
        rw.WriteError(400, err.Error())
        return
    }

    params := r.Form

    // get user info
    users, err = db.GetUsers(params)
    if err != nil {
        rw.WriteError(500, err.Error())
        return
    }

    // form as a JSON response
    m := Meta { Code: 200 }
    resp := JSONResponse { meta: m, data: users }

    // marshall and write
    js, err := json.Marshall(*resp)
    if err != nil {
        rw.WriteError(500, err.Error())
        return
    }

    rw.Write(js)
}

// UsersCreateHandler handles the creation of users.
// Handles PUT and POST to /users.
func (db *sql.DB) UsersCreateHandler(rw http.ResponseWriter, r *http.Request) {
    
    requiredParams := []string{"first_name", 
                               "last_name", 
                               "email", 
                               "access_token", 
                               "verification_token", 
                               "profile_pic_url", 
                               "venmo_id"}

    // parse the form
    if err = r.ParseForm(); err != nil {
        WriteError(400, err.Error())
        return
    }

    params := r.Form

    // verify all params are present
    for _, p := range requiredParams {
        if _, ok := params[p]; !ok {
            rw.WriteError(400, "Missing parameter " + p)
            return
        }
    }

    // create a user
    err := CreateUser(params['first_name'], 
                      params['last_name'], 
                      params['email'], 
                      params['access_token'], 
                      params['verification_token'], 
                      params['profile_pic_url'], 
                      params['venmo_id'])
    if err != nil {
        rw.WriteError(500, "Failed to create user: " + err.Error())
        return
    }

    rw.WriteSuccess()

}

// UserShowHandler handles display of user info by id.
// Handles GET to /user/{id}.
func (db *sql.DB) UserShowHandler(rw http.ResponseWriter, r *http.Request) {

    id := mux.Vars(r)["id"]

    if !db.UserExists(id) {
        rw.WriteError(400, "No user found with id " + strconv.Itoa(id))
        return
    }

    u, err := db.GetUser(id)
    if err != nil {
        rw.WriteError(500, "Unable to retrieve user")
        return
    }

    // form as a JSON response
    m := Meta { Code: 200 }
    resp := JSONResponse { meta: m, data: *u }

    // marshall and write
    js, err := json.Marshall(*resp)
    if err != nil {
        rw.WriteError(500, err.Error())
        return
    }

    rw.Write(js)

}

// UserDeleteHandler handles the deletion of users.
// Handles DELETE at /users/{id}.
func (db *sql.DB) UserDeleteHandler(rw http.ResponseWriter, r *http.Request) {

    id := mux.Vars(r)["id"]

    if !db.UserExists(id) {
        rw.WriteError(400, "No user found with id " + strconv.Itoa(id))
        return
    }

    if err := db.DeleteUser(id); err != nil {
        rw.WriteError(500, "Failed to delete user")
        return
    }

    rw.WriteSuccess()

}

// UserBetsHandler gets the bets a user is participates in.
// Handles GET to /users/{id}/bets.
func (db *sql.DB) UserBetsHandler(rw http.ResponseWriter, r *http.Request) {

    id := mux.Vars(r)["id"]

    if !db.UserExists(id) {
        rw.WriteError(400, "No user found with id " + strconv.Itoa(id))
        return
    }

    bets, err := GetUserBets(id)
    if err != nil {
        rw.WriteError(500, "Failed to get bets for the given user")
        return
    }

    // form as a JSON response
    m := Meta { Code: 200 }
    resp := JSONResponse { meta: m, data: bets }

    // marshall and write
    js, err := json.Marshall(*resp)
    if err != nil {
        rw.WriteError(500, err.Error())
        return
    }

    rw.Write(js)
}

// UserWitnessingHandler gets the bets a user is a witness for.
// Handles GET to /users/{id}/bets.
func (db *sql.DB) UserWitnessingHandler(rw http.ResponseWriter, r *http.Request) {

    id := mux.Vars(r)["id"]

    if !db.UserExists(id) {
        rw.WriteError(400, "No user found with id " + strconv.Itoa(id))
        return
    }

    bets, err := GetUserWitnessing(id)
    if err != nil {
        rw.WriteError(500, "Failed to get bets for the given user")
        return
    }

    // form as a JSON response
    m := Meta { Code: 200 }
    resp := JSONResponse { meta: m, data: bets }

    // marshall and write
    js, err := json.Marshall(*resp)
    if err != nil {
        rw.WriteError(500, err.Error())
        return
    }

    rw.Write(js)
}

// BetsShowHandler handles display of many bets.
// Handles GET to /bets.
func (db *sql.DB) BetsShowHandler(rw http.ResponseWriter, r *http.Request) {

    var bets []Bet 
    var err error

    // parse the form
    if err = r.ParseForm(); err != nil {
        rw.WriteError(400, err.Error())
        return
    }

    params := r.Form

    // get user info
    bets, err = db.GetBets(params)
    if err != nil {
        rw.WriteError(500, err.Error())
        return
    }

    // form as a JSON response
    m := Meta { Code: 200 }
    resp := JSONResponse { meta: m, data: bets }

    // marshall and write
    js, err := json.Marshall(*resp)
    if err != nil {
        rw.WriteError(500, err.Error())
        return
    }

    rw.Write(js)
}

// BetsCreateHandler handles creation of bets.
// Handles PUT and POST to /bets.
// Includes functionality for charging both parties over Venmo.
func (db *sql.DB) BetsCreateHandler(rw http.ResponseWriter, r *http.Request) {

    requiredParams := []string{"bettor_id", 
                               "betted_id", 
                               "witness_id", 
                               "witness_id", 
                               "title", 
                               "desc", 
                               "expire_on",
                               "status",
                               "amount"}

    // parse the form
    if err = r.ParseForm(); err != nil {
        WriteError(400, err.Error())
        return
    }

    params := r.Form

    // verify all params are present
    for _, p := range requiredParams {
        if _, ok := params[p]; !ok {
            rw.WriteError(400, "Missing parameter " + p)
            return
        }
    }

    // create a user
    err := CreateUser(params['bettor_id'], 
                      params['betted_id'], 
                      params['witness_id'], 
                      params['winner_id'], 
                      params['title'], 
                      params['desc'], 
                      params['expire_on'],
                      params['status'],
                      params['amount'])
    if err != nil {
        rw.WriteError(500, "Failed to create bet: " + err.Error())
        return
    }

    rw.WriteSuccess()

}

// BetShowHandler displays info for a bet.
// Handles GET to /bets/{id}.
func (db *sql.DB) BetShowHandler(rw http.ResponseWriter, r *http.Request) {

    id := mux.Vars(r)["id"]

    if !db.BetExists(id) {
        rw.WriteError(400, "No bet found with id " + strconv.Itoa(id))
        return
    }

    b, err := db.GetBet(id)
    if err != nil {
        rw.WriteError(500, "Unable to retrieve bet")
        return
    }

    // form as a JSON response
    m := Meta { Code: 200 }
    resp := JSONResponse { meta: m, data: *b }

    // marshall and write
    js, err := json.Marshall(*resp)
    if err != nil {
        rw.WriteError(500, err.Error())
        return
    }

    rw.Write(js)
}

// BetDeleteHandler handles deletion of bets.
// Handles DELETE to /bets/{id}.
func (db *sql.DB) BetDeleteHandler(rw http.ResponseWriter, r *http.Request) {

    id := mux.Vars(r)["id"]

    if !db.BetExists(id) {
        rw.WriteError(400, "No bet found with id " + strconv.Itoa(id))
        return
    }

    if err := db.DeleteBet(id); err != nil {
        rw.WriteError(500, "Failed to delete bet")
        return
    }

    rw.WriteSuccess()
}

// BetStatusHandler handles changing the status of a bet.
// Handles POST to /bet/{id}/status.
// Must implement:
//  - Pending created automatically on create
//  - Declined allowed by betted and witness
//  - Settled allowed by witness
// Includes requests to Venmo to payout on status = settled.
func (db *sql.DB) BetStatusHandler(rw http.ResponseWriter, r *http.Request) {

    // check id
    id := mux.Vars(r)["id"]

    if !db.BetExists(id) {
        rw.WriteError(400, "No bet found with id " + strconv.Itoa(id))
        return
    }

    // parse the form
    if err = r.ParseForm(); err != nil {
        rw.WriteError(400, err.Error())
        return
    }

    params := r.Form

    // update status
    status, ok := params["status"]
    if !ok {
        rw.WriteError(400, "Required parameter 'status' not provided")
        return
    }

    winnerId := -1
    settled = staus == "settled"

    if settled {
        winnerId, ok = params["winner_id"]
        if !ok {
            rw.WriteError(400, "Required parameter for status settled 'winner_id' not provided")
            return
        }
    }

    err := UpdateBetStatus(status, id, winner_id)
    if err != nil {
        rw.WriteError(500, "Failed to update bet status")
        return
    }

    if settled {
        // Call function to charge bettor_id and betted_id
    }
}

/* Basic Responses */

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