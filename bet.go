package main

import (
    "database/sql"
    "errors"
    "fmt"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

// A Bet represents basic info about a bet.
type Bet struct {
    Id int                `json:"id"`
    BettorId int          `json:"bettor_id"`
    BettedId int          `json:"betted_id"`
    WitnessId int         `json:"witness_id"`
    WinnerId int          `json:"winner_id"`
    Title string          `json:"title"`
    Desc string           `json:"desc"`
    CreatedOn time.Time   `json:"created_on"`
    ExpiresOn time.Time   `json:"expires_on"`
    Status string         `json:"status"`
    Amount int            `json:"amount"` // in cents
}

// Creates a bet.
func (db *sql.DB) CreateBet(bettorId int, 
                            bettedId int, 
                            witnessId int, 
                            winnerId int, 
                            title string, 
                            description string, 
                            expireOn time.Duration, 
                            status string, 
                            amount int) error {

    q := "insert into bet (id, bettor_id, betted_id, witness_id, " + 
         "winner_id, title, description, created,_on, expire_on, " +
         "status, amount values (?)"

    query, err := db.Prepare(q)
    if err != nil {
        return errors.New("Error when preparing the CreateBet query")
    }

    q := fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s, %s, %s", bettorId, 
                                                           bettedId, 
                                                           witnessId, 
                                                           winnerId, 
                                                           title, 
                                                           description, 
                                                           expireOn, 
                                                           status, 
                                                           amount)
    
    _, err := query.Exec(q)
    if err != nil{
        return errors.New("Error when executing the CreateBet query")
    }

    return nil
}

//Retrieves a list of bets by specific query parameters
func (db *sql.DB) RetrieveBets(params map[string] string) ([]Bet, error){
    bets := make([]Bet)
    var b Bet

    query := "select * from bets where "
    for k, v := range params{
        query += fmt.Sprintf("%s = %s" k, v)
        query += "and "
    query += "is_deleted = 0"
    }

    stmt, err := db.Prepare(query)
    if err != nil{
        return errors.New("Error when preparing the RetrieveBets query")
    }
    defer stmt.Close()

    rows, err := stmt.Query()
    if err != nil{
        return errors.New("Error when executing the RetrieveBets query")
    }
    defer rows.Close()

    for rows.Next() {
        err = rows.Scan(&b.Id, 
                        &b.BettorId, 
                        &b.BettedID, 
                        &b.WitnessID, 
                        &b.WinnerID, 
                        &b.Title, 
                        &b.Desc, 
                        &b.CreatedOn, 
                        &b.ExpireOn, 
                        &b.Status, 
                        &b.Amount)
        bets.append(b)
    }
    if err = rows.Err(); err != nil{
        return errors.New("Error when scanning the RetrieveBets query")
    }

    return bets, nil
}

// Retrieves a specific bet by it's id in the database.
func (db *sql.DB) RetrieveBet(id int) (*Bet, error){
    var b Bet
    err := db.QueryRow("select * from bet where id = ?", id).Scan(&b.Id, 
                                                                  &b.BettorId, 
                                                                  &b.BettedID, 
                                                                  &b.WitnessID, 
                                                                  &b.WinnerID, 
                                                                  &b.Title, 
                                                                  &b.Desc, 
                                                                  &b.CreatedOn, 
                                                                  &b.ExpireOn, 
                                                                  &b.Status, 
                                                                  &b.Amount)
    if err != nil{
        return errors.New("Error when executing the RetrieveBet query")
    }

    return &b, nil
}

// Deletes a bet.
// We don't delete the row from the actual table. Instead we toggle the
// is_deleted attribute in the database.
func (db *sql.DB) DeleteBet(id int) error{

    query, err := db.Prepare("update bets set is_deleted = 1 where id = ?")
    if err != nil{
        return errors.New("Error when preparing the DeleteBet query")
    }

    _, err := query.Exec(id)
    if err != nil{
        return errors.New("Error when executing the DeleteBet query")
    }

    return nil
}

// Updates a specific bet with it's winner.
func (db *sql.DB) UpdateBetWinner(id int, winnerId int) error {

    query, err := db.Prepare("update bets set winner_id = ? where id = ? and is_deleted = 0")
    if err != nil {
        return errors.New("Error when preparing the UpdateBetWinner query")
    }

    _, err := query.Exec(winnerId, id)
    if err != nil{
        return errors.New("Error when executing the UpdateBetWinner query")
    }

    return nil
}
