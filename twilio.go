package main

import (
    "errors"
    "fmt"
    "net/http"
    "net/url"
    "os"
    "strings"

    _ "github.com/go-sql-driver/mysql"
)

// SendVerificationMsg sends a text message with a user's verification token so they can confirm their phone number.
func (db *MyDB) SendVerificationMsg(accessToken string, phoneNumber string) error {

    verificationToken, err := db.GetVerificationTokenFromAccessToken(accessToken)
    if err != nil {
        return err
    }

    msg := fmt.Sprintf("Your Bettor verification id is: %s", verificationToken)
    SendTwilioMsg(phoneNumber, msg)

    return nil
}

// SendTwilioMsg sends a text message from the Twilio API.
func SendTwilioMsg(phoneNumber string, message string) error{

    accountSid := "AC4b7b097d333a0d6490fff5d1098db453"
    authToken := os.Getenv("TWILIO_SECRET_KEY")
    urlString := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)

    url_params := url.Values{}
    url_params.Set("From", "+19782212680")
    url_params.Set("To", fmt.Sprintf("+1%s", phoneNumber))
    url_params.Set("Body", message)
    req_body := *strings.NewReader(url_params.Encode())

    client := &http.Client{}
    req, err := http.NewRequest("POST", urlString, &req_body)
    if err != nil{
        return errors.New("Unable to create NewRequest: " + err.Error())
    }
    req.SetBasicAuth(accountSid, authToken)
    req.Header.Add("Accept", "application/json")
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    _, err = client.Do(req)
    if err != nil{
        return errors.New("Verification text message failed to send: " + err.Error())
    }

    return nil
}

// GetVerificationTokenFromAccessToken returns the verification token based on a 
// given user's Venmo access token.
func (db *MyDB) GetVerificationTokenFromAccessToken(accessToken string) (string, error){
    var verificationToken string
    err := db.QueryRow("select verification_token from users where access_token = ?", accessToken).Scan(&verificationToken)
    if err != nil{
        return "", errors.New("Failed while querying for the verification token: " + err.Error())
    }

    return verificationToken, nil
}