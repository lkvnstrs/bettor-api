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
    CreatedOn time.Time   `json:"created_at"`
    ExpireOn time.Time    `json:"expire_at"`
    Status string         `json:"status"`
    Amount int            `json:"amount"` // in cents
}

// Creates a bet.
func (db *sql.DB) CreateBet(bettor_id int, betted_id int, witness_id int, winner_id int, title string, description string, expire_on time.Duration, status string, amount int) error {
    query, err := db.Prepare("insert into bet (id, bettor_id, betted_id, witness_id, winner_id, title, description, created,_on, expire_on, status, amount values (?)")
    if err != nil {
        return errors.New("Error when preparing the CreateBet query")
    }

    result, err := query.Exec(fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s, %s, %s", bettor_id, betted_id, witness_id, winner_id, title, description, expire_on, status, amount))
    if err != nil{
        return errors.New("Error when executing the CreateBet query")
    }

    lastID, err := result.LastInsertId()
    if err != nil{
        return errors.New("Error retrieving last inserted id for CreateBet query")
    }

    rowCount, err := result.RowsAffected()
    if err != nil{
        return errors.New("Error retrieving total affected rows for CreateBet query")
    }
}

//Retrieves a list of bets by specific query parameters
func (db *sql.DB) RetrieveBets(query_params map[string] string) ([]Bet, error){
    bets := make([]Bet)
    var b Bet

    query := "select * from bets where "
    for k, v := range query_params{
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
        err = rows.Scan(&b.Id, &b.BettorId, &b.BettedID, &b.WitnessID, &b.WinnerID, &b.Title, &b.Desc, &b.CreatedOn, &b.ExpireOn, &b.Status, &b.Amount)
        bets.append(b)
    }
    if err = rows.Err(); err != nil{
        return errors.New("Error when scanning the RetrieveBets query")
    }

    return bets, nil
}

// Retrieves a specific bet by it's id in the database.
func (db *sql.DB) RetrieveBet(id int) (Bet, error){
    var b Bet
    if err := db.QueryRow("select * from bet where id = ?", id).Scan(&b.Id, &b.BettorId, &b.BettedID, &b.WitnessID, &b.WinnerID, &b.Title, &b.Desc, &b.CreatedOn, &b.ExpireOn, &b.Status, &b.Amount); err != nil{
        return errors.New("Error when executing the RetrieveBet query")
    }

    return b, nil
}

// Deletes a bet.
// We don't delete the row from the actual table. Instead we toggle the
// is_deleted attribute in the database.
func (db *sql.DB) DeleteBet(id int) error{
    query, err := db.Prepare("update bets set is_deleted = 1 where id = ?")
    if err != nil{
        return errors.New("Error when preparing the DeleteBet query")
    }

    result, err := query.Exec(id)
    if err != nil{
        return errors.New("Error when executing the DeleteBet query")
    }

    lastID, err := result.LastInsertId()
    if err != nil{
        return errors.New("Error retrieving last inserted id for DeleteBet query")
    }

    rowCount, err := result.RowsAffected()
    if err != nil{
        return errors.New("Error retrieving total affected rows for DeleteBet query")
    }
}

// Updates a specific bet with it's winner.
func (db *sql.DB) UpdateBetWinner(id int, winner_id int) error{
    query, err := db.Prepare("update bets set winner_id = ? where id = ? and is_deleted = 0")
    if err != nil {
        return errors.New("Error when preparing the UpdateBetWinner query")
    }

    result, err := query.Exec(winner_id, id)
    if err != nil{
        return errors.New("Error when executing the UpdateBetWinner query")
    }

    lastID, err := result.LastInsertId()
    if err != nil{
        return errors.New("Error retrieving last inserted id for UpdateBetWinner query")
    }

    rowCount, err := result.RowsAffected()
    if err != nil{
        return errors.New("Error retrieving total affected rows for UpdateBetWinner query")
    }
}
