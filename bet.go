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
    Status string         `json:"status"`
    Amount int            `json:"amount"` // in cents
}

// CreateBet creates a bet.
func (db *MyDB) CreateBet(bettorId int,
                            bettedId int,
                            witnessId int,
                            winnerId int,
                            title string,
                            description string,
                            status string,
                            amount int) error {

    q := "insert into bet (id, bettor_id, betted_id, witness_id, " +
         "winner_id, title, description, created,_on, expire_on, " +
         "status, amount values (?)"

    query, err := db.Prepare(q)
    if err != nil {
        return errors.New("Error when preparing the CreateBet query")
    }

    values := fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s, %s, %s", bettorId,
                                                           bettedId,
                                                           witnessId,
                                                           winnerId,
                                                           title,
                                                           description,
                                                           status,
                                                           amount)

    _, err = query.Exec(values)
    if err != nil{
        return errors.New("Error when executing the CreateBet query")
    }

    return nil
}

// GetBets retrieves a list of bets by specific query parameters.
func (db *MyDB) GetBets(params map[string] string) ([]Bet, error){
    bets := make([]Bet, 0)
    var b Bet

    query := "select * from bets where is_deleted = 0 and "
    for k, v := range params {
        query += fmt.Sprintf("%s = '%s' and ", k, v)
    }

    query = query[:len(query) - 5]

    stmt, err := db.Prepare(query)
    if err != nil{
        return nil, errors.New("Error when preparing the bet query")
    }
    defer stmt.Close()

    rows, err := stmt.Query()
    if err != nil{
        return nil, errors.New("Error when executing the bet query")
    }
    defer rows.Close()

    for rows.Next() {
        err = rows.Scan(&b.Id,
                        &b.BettorId,
                        &b.BettedId,
                        &b.WitnessId,
                        &b.WinnerId,
                        &b.Title,
                        &b.Desc,
                        &b.CreatedOn,
                        &b.Status,
                        &b.Amount)

        bets = append(bets, b)
    }
    if err = rows.Err(); err != nil{
        return nil, errors.New("Error when scanning the RetrieveBets query")
    }

    return bets[0:], nil
}

// GetBet retrieves a specific bet by it's id in the database.
func (db *MyDB) GetBet(id int) (*Bet, error){
    var b Bet
    err := db.QueryRow("select * from bet where id = ?", id).Scan(&b.Id,
                                                                  &b.BettorId,
                                                                  &b.BettedId,
                                                                  &b.WitnessId,
                                                                  &b.WinnerId,
                                                                  &b.Title,
                                                                  &b.Desc,
                                                                  &b.CreatedOn,
                                                                  &b.Status,
                                                                  &b.Amount)
    if err != nil{
        return nil, errors.New("Error when executing the RetrieveBet query")
    }

    return &b, nil
}

// DeleteBet deletes a bet.
// Doesn't delete the row from the actual table. 
// Toggles the is_deleted attribute in the database.
func (db *MyDB) DeleteBet(id int) error{

    query, err := db.Prepare("update bets set is_deleted = 1 where id = ?")
    if err != nil{
        return errors.New("Error when preparing the DeleteBet query")
    }

    _, err = query.Exec(id)
    if err != nil{
        return errors.New("Error when executing the DeleteBet query")
    }

    return nil
}

// UpdateBetStatus updates the status of a bet.
// Status can only set to "declined", "pending", "active", and "settled".
// Status "settled" charges
func (db *MyDB) UpdateBetStatus(id int, status string, winnerId int) error {

    settled := status == "settled"
    q := "update bets set status=?"

    if settled {
        q += ", winner_id=?"
    }

    q += " where id = ? and is_deleted = 0"

    query, err := db.Prepare(q)
    if err != nil {
        return errors.New("Failed to update bet status")
    }

    if settled {
        _, err := query.Exec(status, winnerId, id)
        if err != nil {
            return errors.New("Failed to update bet status")
        }
    } else {
        _, err := query.Exec(status, id)
        if err != nil {
            return errors.New("Failed to update bet status")
        }
    }
    if err != nil{
        return errors.New("Failed to update bet status")
    }

    // if settled {
    //     // 
    // }

    return nil

}

// BetExists checks if a bet with the given id exists.
func (db *MyDB) BetExists(id int) bool {
    var tmp int
    err := db.QueryRow("select id from bets where id = ? and is_deleted = 0", id).Scan(&tmp)
    return err == sql.ErrNoRows
}