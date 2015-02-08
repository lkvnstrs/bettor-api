package main

import (
    "errors"
    "strconv"

    _ "github.com/go-sql-driver/mysql"
)

// A Contact represents a single contact.
type Contact struct {
    DisplayName string  `json:"display_name"`
    Emails []string     `json:"emails"`
    Phones []string     `json:"phones"`
}

// A ContactPair connects a phone number with an associated user id.
type ContactPair struct{
    PhoneNumber string  `json:"phone_number"`
    UserId int          `json:"user_id"`
}

func (db *MyDB) CheckPhoneNumbers(phoneNumbers []string) ([]ContactPair, error){

    contactpairs := make([]ContactPair, 0)
    var cp ContactPair

    query := "select id, phone_number from users where phone_number in ?"

    for phone_number := range phoneNumbers{
        query += strconv.Itoa(phone_number)
    }
    query = query[:len(query) - 2]

    rows, err := db.Query(query)
    if err != nil{
        return nil, errors.New("Error when executing the CheckPhoneNumbers query")
    }
    defer rows.Close()

    for rows.Next(){

        err = rows.Scan(&cp.PhoneNumber, &cp.UserId)
        contactpairs = append(contactpairs, cp)
    }

    if err = rows.Err(); err != nil{
        return nil, errors.New("Error when scanning the CheckPhoneNumbers query")
    }

    return contactpairs[0:], nil
}