package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sillypears/condor-standings/src/log"
	"github.com/sillypears/condor-standings/src/models"

	"regexp"
)

func httpAPI(c *gin.Context) {
	// Local variables
	w := c.Writer

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"Documentation\": \"API Documentation\" }"))

}

func httpEventAPI(c *gin.Context) {
	// Local variables
	w := c.Writer
	var foundEvent models.Event
	// Set the header
	w.Header().Set("Content-Type", "application/json")

	// Parse the event name from the URL
	event := c.Params.ByName("event")
	if event == "" {
		w.Write([]byte("{\"Error\": \"No event found\"}"))
		return
	}
	groups, err := regexp.MatchString("ndwc*", event)
	if err != nil {
		log.Error("Something funky  happened with regex")
	}
	log.Info(groups)
	if groups == true {
		foundEvent, _ = db.EventAPI.GetEventInfoGroups(event)
	} else {
		foundEvent, _ = db.EventAPI.GetEventInfo(event)
	}
	if err != nil {
		log.Error("Couldn't get event info: ", err)
	}

	jsonData, err := json.MarshalIndent(foundEvent, "", "\t")
	if err != nil {
		log.Error("Couldn't generate JSON")
		//w.Write([]byte("Please search for a user"))
		return
	}
	w.Write(jsonData)

}