package main

import (
    "encoding/json"
    "io/ioutil"
    // "log"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
)

/* Handlers */

// ContactsHandler handles parsing contacts for existing users.
// Handles PUT and POST to /contacts.
func (db *MyDB) ContactsHandler(rw http.ResponseWriter, r *http.Request) {

    var contacts []Contact
    var contactpairs []ContactPair
    var err error

    // parse the data
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&contacts); err != nil {
        WriteError(rw, 500, err.Error())
        return
    }

    phonenumbers := make([]string, len(contacts))
    i := 0
    for _, c := range contacts {
        for _, phone := range c.Phones {
            phonenumbers[i] = phone
            i++
        }
    }

    contactpairs, err = db.CheckPhoneNumbers(phonenumbers[0:])
    if err != nil {
        WriteError(rw, 500, err.Error())
        return
    }

    // form as a JSON response
    m := M{ Code: 200 }
    resp := JSONResponse { Meta: m, Data: contactpairs }

    // marshall and write
    js, err := json.Marshal(resp)
    if err != nil {
        WriteError(rw, 500, err.Error())
        return
    }

    rw.WriteHeader(200)
    rw.Header().Set("Access-Control-Allow-Origin", "*")
    rw.Write(js)
}

// VerificationHandler handles the verification of a user's phone number.
// Handles POST to /verify.
func (db *MyDB) VerificationHandler(rw http.ResponseWriter, r *http.Request) {

    // parse the data
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        WriteError(rw, 500, "Failed to parse body: " + err.Error())
        return 
    }

    var params map[string]string
    err = json.Unmarshal(body, &params)

    // Get access token and verification token
    accessToken, ok := params["access_token"]
    if !ok {
        WriteError(rw, 400, "Parameter 'access_token' is required to verify")
        return
    }

    verificationToken, ok := params["verification_token"]
    if !ok {
        WriteError(rw, 400, "Parameter 'verification_token' is required to verify")
        return
    }

    // Verify user
    if err := db.VerifyUser(accessToken, verificationToken); err != nil {
        WriteError(rw, 400, err.Error())
        return
    }

    WriteSuccess(rw)

}

// UsersShowHandler handles display of users.
// Handles GET to /users.
func (db *MyDB) UsersShowHandler(rw http.ResponseWriter, r *http.Request) {

    var users []User 

    // parse the form
    if err := r.ParseForm(); err != nil {
        WriteError(rw, 400, err.Error())
        return
    }

    // get user info
    params := make(map[string]string)
    for k, v := range r.Form {
        params[k] = v[0]
    }

    users, err := db.GetUsers(params)
    if err != nil {
        WriteError(rw, 500, err.Error())
        return
    }

    // form as a JSON response
    m := M{ Code: 200 }
    resp := JSONResponse { Meta: m, Data: users }

    // marshall and write
    js, err := json.Marshal(resp)
    if err != nil {
        WriteError(rw, 500, err.Error())
        return
    }

    rw.WriteHeader(200)
    rw.Header().Set("Access-Control-Allow-Origin", "*")
    rw.Write(js)
}

// UsersCreateHandler handles the creation of users.
// Handles PUT and POST to /users.
func (db *MyDB) UsersCreateHandler(rw http.ResponseWriter, r *http.Request) {

    // parse the data
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        WriteError(rw, 500, "Failed to parse body: " + err.Error())
        return
    }

    var params map[string]string
    err = json.Unmarshal(body, &params)

    // Get access token and phone number
    accessToken, ok := params["access_token"]
    if !ok {
        WriteError(rw, 400, "Required parameter 'access_token' is not present")
        return
    }

    phoneNumber, ok := params["phone_number"]
    if !ok {
        WriteError(rw, 400, "Required parameter 'phone_number' is not present")
        return
    }

    // // request user info from venmo
    info, err := GetVenmoInfo(accessToken)
    if err != nil {
        WriteError(rw, 500, "Failed request for info from Venmo: " + err.Error())
        return
    }

    // create a user
    err = db.CreateUser(info["first_name"], 
                        info["last_name"],
                        info["email"], 
                        accessToken,
                        info["profile_pic_url"], 
                        info["venmo_id"],
                        phoneNumber)
    if err != nil {
        WriteError(rw, 500, "Failed to create user: " + err.Error())
        return
    }

    // send twilio
    err = db.SendVerificationMsg(accessToken, phoneNumber)
    if err != nil {
        WriteError(rw, 500, "Failed to send Twilio message: " + err.Error())
        return
    }

    WriteSuccess(rw)

}

// UserShowHandler handles display of user info by id.
// Handles GET to /user/{id}.
func (db *MyDB) UserShowHandler(rw http.ResponseWriter, r *http.Request) {

    id, _ := strconv.Atoi(mux.Vars(r)["id"])

    if !db.UserExists(id) {
        WriteError(rw, 400, "No user found with id " + strconv.Itoa(id))
        return
    }

    u, err := db.GetUser(id)
    if err != nil {
        WriteError(rw, 500, "Unable to retrieve user: " + err.Error())
        return
    }

    // form as a JSON response
    m := M{ Code: 200 }
    resp := JSONResponse { Meta: m, Data: *u }

    // marshall and write
    js, err := json.Marshal(resp)
    if err != nil {
        WriteError(rw, 500, err.Error())
        return
    }

    rw.WriteHeader(200)
    rw.Header().Set("Access-Control-Allow-Origin", "*")
    rw.Write(js)

}

// UserUpdateHandler handles updating a user.
// Handles POST at /users/{id}.
func (db *MyDB) UserUpdateHandler(rw http.ResponseWriter, r *http.Request) {

    id, _ := strconv.Atoi(mux.Vars(r)["id"])

    if !db.UserExists(id) {
        WriteError(rw, 400, "No user found with id " + strconv.Itoa(id))
        return
    }

    // parse the data
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        WriteError(rw, 500, "Failed to parse body: " + err.Error())
        return 
    }

    var params map[string]string
    err = json.Unmarshal(body, &params)

    err = db.UpdateUser(id, params); if err != nil {
        WriteError(rw, 400, err.Error())
        return
    }

    WriteSuccess(rw)
}

// UserDeleteHandler handles the deletion of users.
// Handles DELETE at /users/{id}.
func (db *MyDB) UserDeleteHandler(rw http.ResponseWriter, r *http.Request) {

    id, _ := strconv.Atoi(mux.Vars(r)["id"])

    if !db.UserExists(id) {
        WriteError(rw, 400, "No user found with id " + strconv.Itoa(id))
        return
    }

    if err := db.DeleteUser(id); err != nil {
        WriteError(rw, 500, "Failed to delete user")
        return
    }

    WriteSuccess(rw)

}

// UserBetsHandler gets the bets a user is participates in.
// Handles GET to /users/{id}/bets.
func (db *MyDB) UserBetsHandler(rw http.ResponseWriter, r *http.Request) {

    id, _ := strconv.Atoi(mux.Vars(r)["id"])

    if !db.UserExists(id) {
        WriteError(rw, 400, "No user found with id " + strconv.Itoa(id))
        return
    }

    bets, err := db.GetUserBets(id)
    if err != nil {
        WriteError(rw, 500, "Failed to get bets for the given user")
        return
    }

    // form as a JSON response
    m := M{ Code: 200 }
    resp := JSONResponse { Meta: m, Data: bets }

    // marshall and write
    js, err := json.Marshal(resp)
    if err != nil {
        WriteError(rw, 500, err.Error())
        return
    }

    rw.WriteHeader(200)
    rw.Header().Set("Access-Control-Allow-Origin", "*")
    rw.Write(js)
}

// UserWitnessingHandler gets the bets a user is a witness for.
// Handles GET to /users/{id}/bets.
func (db *MyDB) UserWitnessingHandler(rw http.ResponseWriter, r *http.Request) {

    id, _ := strconv.Atoi(mux.Vars(r)["id"])

    if !db.UserExists(id) {
        WriteError(rw, 400, "No user found with id " + strconv.Itoa(id))
        return
    }

    bets, err := db.GetUserWitnessing(id)
    if err != nil {
        WriteError(rw, 500, "Failed to get bets for the given user")
        return
    }

    // form as a JSON response
    m := M{ Code: 200 }
    resp := JSONResponse { Meta: m, Data: bets }

    // marshall and write
    js, err := json.Marshal(resp)
    if err != nil {
        WriteError(rw, 500, err.Error())
        return
    }

    rw.WriteHeader(200)
    rw.Header().Set("Access-Control-Allow-Origin", "*")
    rw.Write(js)
}

// BetsShowHandler handles display of many bets.
// Handles GET to /bets.
func (db *MyDB) BetsShowHandler(rw http.ResponseWriter, r *http.Request) {

    var bets []Bet 
    var err error

    // parse the form
    if err = r.ParseForm(); err != nil {
        WriteError(rw, 400, err.Error())
        return
    }

    params := make(map[string]string)
    for k, v := range r.Form {
        params[k] = v[0]
    }

    // get user info
    bets, err = db.GetBets(params)
    if err != nil {
        WriteError(rw, 500, err.Error())
        return
    }

    // form as a JSON response
    m := M{ Code: 200 }
    resp := JSONResponse { Meta: m, Data: bets }

    // marshall and write
    js, err := json.Marshal(resp)
    if err != nil {
        WriteError(rw, 500, err.Error())
        return
    }

    rw.WriteHeader(200)
    rw.Header().Set("Access-Control-Allow-Origin", "*")
    rw.Write(js)
}

// BetsCreateHandler handles creation of bets.
// Handles PUT and POST to /bets.
// Includes functionality for charging both parties over Venmo.
func (db *MyDB) BetsCreateHandler(rw http.ResponseWriter, r *http.Request) {

    requiredParams := []string{"access_token", 
                               "betted_id", 
                               "witness_id",
                               "title",
                               "amount"}

    // parse the data
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        WriteError(rw, 500, "Failed to parse body: " + err.Error())
        return 
    }

    var params map[string]string
    err = json.Unmarshal(body, &params)

    // verify required params
    for _, p := range requiredParams {
        if _, ok := params[p]; !ok {
            WriteError(rw, 400, "Missing parameter " + p)
            return
        }
    }

    // request sets
    bettorId, err := db.GetIdByAccessToken(params["access_token"])
    if err != nil {
        WriteError(rw, 400, err.Error())
        return
    }

    bettedId, _ := strconv.Atoi(params["betted_id"])
    witnessId, _ := strconv.Atoi(params["witness_id"])
    amount, _ := strconv.Atoi(params["amount"])
    title := params["title"]
   
    // defaults to
    winnerId := 0
    desc := ""
    status := "pending"

    // create a user
    err = db.CreateBet(bettorId, 
                       bettedId, 
                       witnessId, 
                       winnerId,
                       title, 
                       desc, 
                       status,
                       amount)
    if err != nil {
        WriteError(rw, 500, "Failed to create bet: " + err.Error())
        return
    }

    WriteSuccess(rw)

}

// BetShowHandler displays info for a bet.
// Handles GET to /bets/{id}.
func (db *MyDB) BetShowHandler(rw http.ResponseWriter, r *http.Request) {

    id, _ := strconv.Atoi(mux.Vars(r)["id"])

    if !db.BetExists(id) {
        WriteError(rw, 400, "No bet found with id " + strconv.Itoa(id))
        return
    }

    b, err := db.GetBet(id)
    if err != nil {
        WriteError(rw, 500, "Unable to retrieve bet")
        return
    }

    // form as a JSON response
    m := M{ Code: 200 }
    resp := JSONResponse { Meta: m, Data: *b }

    // marshall and write
    js, err := json.Marshal(resp)
    if err != nil {
        WriteError(rw, 500, err.Error())
        return
    }

    rw.WriteHeader(200)
    rw.Header().Set("Access-Control-Allow-Origin", "*")
    rw.Write(js)
}

// BetDeleteHandler handles deletion of bets.
// Handles DELETE to /bets/{id}.
func (db *MyDB) BetDeleteHandler(rw http.ResponseWriter, r *http.Request) {

    id, _ := strconv.Atoi(mux.Vars(r)["id"])

    if !db.BetExists(id) {
        WriteError(rw, 400, "No bet found with id " + strconv.Itoa(id))
        return
    }

    if err := db.DeleteBet(id); err != nil {
        WriteError(rw, 500, "Failed to delete bet")
        return
    }

    WriteSuccess(rw)
}

// BetStatusHandler handles changing the status of a bet.
// Handles POST to /bet/{id}/status.
// Must implement:
//  - Pending created automatically on create
//  - Declined allowed by betted and witness
//  - Settled allowed by witness
// Includes requests to Venmo to payout on status = settled.
func (db *MyDB) BetStatusHandler(rw http.ResponseWriter, r *http.Request) {

    var settled bool
    var err error

    // check id
    id, _ := strconv.Atoi(mux.Vars(r)["id"])

    if !db.BetExists(id) {
        WriteError(rw, 400, "No bet found with id " + strconv.Itoa(id))
        return
    }

    // parse the data
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        WriteError(rw, 500, "Failed to parse body: " + err.Error())
        return 
    }

    var params map[string]string
    err = json.Unmarshal(body, &params)

    // update status
    _, ok := params["status"]
    if !ok {
        WriteError(rw, 400, "Required parameter 'status' not provided")
        return
    }

    status := params["status"]

    winnerId := -1
    settled = (status == "settled")

    if settled {
        _, ok = params["winner_id"]
        if !ok {
            WriteError(rw, 400, "Required parameter for status settled 'winner_id' not provided")
            return
        }
        winnerId, err = strconv.Atoi(params["winner_id"])
        if err != nil {
            WriteError(rw, 400, "Parameter 'winner_id' must be an integer")
            return
        }
    }

    err = db.UpdateBetStatus(id, status, winnerId)
    if err != nil {
        WriteError(rw, 500, "Failed to update bet status")
        return
    }

    // if settled {
    //     // Call function to charge bettor_id and betted_id
    // }

    WriteSuccess(rw)
}


// BetsHookHandler handles requests from the Venmo webhook.
// Handles PUT and POST to /bets/hook.
func (db *MyDB) BetsHookHandler(rw http.ResponseWriter, r *http.Request) {

    rw.WriteHeader(200)
    rw.Write([]byte("Not implemented yet"))
}

/* Basic Responses */

// WriteError writes a JSON-formatted error response to a ResponseWriter.
func WriteError(rw http.ResponseWriter, code int, errMsg string) {
    js, _ := json.Marshal(*GenerateError(code, errMsg))

    rw.Header().Set("Access-Control-Allow-Origin", "*")
    rw.WriteHeader(code)
    rw.Write(js)
}

// WriteSuccess writes a JSON-formatted success response to a ResponseWriter.
func WriteSuccess(rw http.ResponseWriter) {
    code := 200
    js, _ := json.Marshal(JSONResponse{ Meta: M{ Code: code, ErrorMessage: "" }})

    rw.Header().Set("Access-Control-Allow-Origin", "*")
    rw.WriteHeader(code)
    rw.Write(js)
}