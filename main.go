package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database, _ := sql.Open("sqlite3", "./pulse.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS activities (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT, timestamp INTEGER, elapsed_time INTEGER, type INTEGER, distance INTEGER, name STRING)")
	statement.Exec()

	url := "https://www.strava.com/api/v3/clubs/pulselive/activities"
	token := "nope"

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + token

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error on response.\n[ERROR] -", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	type stravaAthlete struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
	}

	type stravaActivity struct {
		Athlete      stravaAthlete `json:"athlete"`
		Name         string        `json:"name"`
		Distance     float64       `json:"distance"`
		ElapsedTime  int           `json:"elapsed_time"`
		ActivityType string        `json:"type"`
	}

	jsonBody := []stravaActivity{}
	unmarshalErr := json.Unmarshal([]byte(body), &jsonBody)
	if unmarshalErr != nil {
		fmt.Println(err)
	}

	statement, _ = database.Prepare("INSERT INTO activities (firstname, lastname, timestamp, elapsed_time, type, distance, name) VALUES (?, ?, ?, ?, ?, ?, ?)")

	// Loop over list in reverse
	for i := len(jsonBody) - 1; i >= 0; i-- {
		statement.Exec(jsonBody[i].Athlete.FirstName, jsonBody[i].Athlete.LastName, time.Now().Unix(), jsonBody[i].ElapsedTime, jsonBody[i].ActivityType, jsonBody[i].Distance, jsonBody[i].Name)
	}

	// Do more http requests if we don't have the oldest entry
	// Insert if not present from oldest to newest
}
