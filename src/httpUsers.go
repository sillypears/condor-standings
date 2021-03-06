package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"strings"
	"fmt"
	"math" 
	"github.com/gin-gonic/gin"
	"github.com/sillypears/condor-standings/src/log"
	"github.com/sillypears/condor-standings/src/models"
)

//  --------------------------------
//  The landing page for the website
//  --------------------------------

func formatTime(t int) (string) {
	w := float64(t)/float64(100)
	ms := int((float64(t/100) - w) * -100)
	s := ((t/100) % 60)
	m := ((t/(100*60)) % 60)
	h := ((t/(100*60*60)) % 24)
	var fTime string
	if h > 0 {
		fTime += fmt.Sprintf("%02d:", h)
	}
	fTime += fmt.Sprintf("%02d:%02d.%02d", m, s, ms)
	return fTime
}

func httpUser(c *gin.Context) {
	w := c.Writer
	
	data := TemplateData{
		Title: 			 "CoNDUIT 38 Users",
		
	}
	httpServeTemplate(w, "users", data)

}



func httpUserInfo(c *gin.Context) {
	w := c.Writer

	twitchUsername := c.Params.ByName("user")
	eventName:= ""
	twitterUsername := ""
	if twitchUsername == "" {
		w.Write([]byte("{\"Error\": \"No username found\"}"))
		return
	}

	url := "https://some.pizza/api/event"
	urlClient := http.Client{
		Timeout: time.Second * 2,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error(err)
	}
	req.Header.Set("User-Agent", "it's me!")
	res, getErr := urlClient.Do(req)
	if getErr != nil {
		log.Error(getErr)
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Error(readErr)
	}
	var EventTables []models.ReturnedTable
	jsonErr := json.Unmarshal(body, &EventTables)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
	eventName = EventTables[0].EventName

	url = fmt.Sprintf("https://some.pizza/api/user/%s/%s", eventName, twitchUsername)
	log.Info(url)
	urlClient = http.Client{
		Timeout: time.Second * 2,
	}
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error(err)
	}
	req.Header.Set("User-Agent", "it's me!")
	res, getErr = urlClient.Do(req)
	if getErr != nil {
		log.Error(getErr)
	}
	body, readErr = ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Error(readErr)
	}
	
	// Build the results struct
	var UserMatches []models.Match
	jsonErr = json.Unmarshal(body, &UserMatches)
	if jsonErr != nil {
		log.Error(jsonErr)
	}

	userWins := 0
	userLosses := 0
	addTime := 0
	fstTime := 0
	avgTime := 0
	userWinsPerc := 0.00
	userLossPerc := 0.00
	for _, race := range UserMatches {
		winner := race.RaceWinner

		if strings.ToLower(twitchUsername) == strings.ToLower(race.Racer1Name) && winner == 1 {
			userWins++
			addTime += race.RaceTime
			if fstTime > race.RaceTime || fstTime == 0  {
				fstTime = race.RaceTime
			}
		} else if strings.ToLower(twitchUsername) == strings.ToLower(race.Racer2Name) && winner == 2 {
			userWins++
			addTime += race.RaceTime
			if fstTime > race.RaceTime  || fstTime == 0 {
				fstTime = race.RaceTime
			}
		} else {
			userLosses++
		}

	}

	if userWins > 0 {

		avgTime = addTime/userWins
		userWinsPerc = math.Round( (float64(userWins) / float64(len(UserMatches)) * 10000 ) / 100 )
		userLossPerc = math.Round( (float64(userLosses) / float64(len(UserMatches)) * 10000 ) / 100 )
	}
	
	avgTimeF := formatTime(avgTime)
	fstTimeF := formatTime(fstTime)

	data := TemplateData{
		Title: 			 "User Info",
		TwitchUsername:	 twitchUsername,
		TwitterUsername: twitterUsername,
		UserMatchInfo:   UserMatches,
		UserWins:        userWins,
		WinPerc:         userWinsPerc,
		LossPerc:        userLossPerc,
		UserLosses:      userLosses,
		AvgTime:		 avgTime,
		AvgTimeF:        avgTimeF,
		FastTime:        fstTime,
		FastTimeF:       fstTimeF,
		FoundTables:	 EventTables,
	}

	httpServeTemplate(w, "user", data)
}
