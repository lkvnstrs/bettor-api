package main

import (
    "database/sql"
    "errors"
    "time"

    _ "github.com/lib/pq"
)

// A Userinfo represents basic info about a user.
type Bet struct {
    Id string               `json:"id"`
    BettorId string         `json:"bettor_id"`
    BettedId string         `json:"betted_id"`
    WitnessId string        `json:"witness_id"`
    WinnerId string         `json:"winner_id"`
    Title string            `json:"title"`
    Desc string             `json:"desc"`
    CreatedAt time.Duration `json:"created_at"`
    ExpireAt time.Duration  `json:"expire_at"`
    Status string           `json:"status"`
    Amount int              `json:"amount"` // in cents
}