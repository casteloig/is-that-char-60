package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	sentry "github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
)

var auth_uri = "https://us.battle.net/oauth/token"

type DataAuth struct {
	AccessToken string `json:"access_token"`
}

type DataChar struct {
	Level int `json:"level"`
}

func auth() string {

	client := &http.Client{}
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	req, err := http.NewRequest("POST", auth_uri, strings.NewReader(form.Encode()))
	if err != nil {
		sentry.CaptureMessage("Error in NewRequest")
		log.WithFields(log.Fields{
			"func": "auth",
		}).Error("Error in NewRequest")
		fmt.Println("Error in NewRequest")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(os.Getenv("CLIENT"), os.Getenv("SECRET"))

	resp, err := client.Do(req)
	if err != nil {
		sentry.CaptureMessage("Error in clientDo")
		log.WithFields(log.Fields{
			"func": "auth",
		}).Error("Error in clientDo")
		fmt.Println("Error in clientDo")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		sentry.CaptureMessage("Error in ReadAll")
		log.WithFields(log.Fields{
			"func": "auth",
		}).Error("Error while reading body response")
		fmt.Println("Error in ReadAll")
	}

	dataAuth1 := DataAuth{}
	_ = json.Unmarshal(body, &dataAuth1)

	return dataAuth1.AccessToken

}

func isChar60(realm string, charName string, region string, bearerToken string) (int, int) {

	uri_char := "https://" + strings.ToLower(region) + ".api.blizzard.com/profile/wow/character/" +
		strings.ToLower(realm) + "/" + strings.ToLower(charName) + "?namespace=profile-" + strings.ToLower(region) + "&locale=en_" + strings.ToUpper(region)

	var responseCode int

	client := &http.Client{}
	req, err := http.NewRequest("GET", uri_char, nil)
	if err != nil {
		sentry.CaptureMessage("Error in NewRequest")
		log.WithFields(log.Fields{
			"func": "isChar60",
		}).Error("Error NewRequest")
		fmt.Println("Error in NewRequest")
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, err := client.Do(req)
	if err != nil {
		sentry.CaptureMessage("Error in clientDo")
		log.WithFields(log.Fields{
			"func": "isChar60",
		}).Error("Error in clientDo")
		fmt.Println("Error in clientDo")
	}

	responseCode = resp.StatusCode
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		sentry.CaptureMessage("Error in ReadAll")
		log.WithFields(log.Fields{
			"func": "isChar60",
		}).Error("Error while reading body response")
		fmt.Println("Error in ReadAll")
	}

	dataChar1 := DataChar{}
	_ = json.Unmarshal(body, &dataChar1)

	return dataChar1.Level, responseCode
}
