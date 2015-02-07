package main

import (
    "database/sql"
    "errors"

    _ "github.com/lib/pq"
)

// A User represents basic info about a user.
type User struct {
    Id string               `json:"id"`
    FirstName string        `json:"first_name"`
    LastName string         `json:"last_name"`
    Email string            `json:"email"`
    AccessToken string      `json:"access_token"`
    ProfilePicUrl string    `json:"profile_pic_url"`
    VenmoId string          `json:"venmo_id"`
}