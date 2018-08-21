// @SubApi API [/api]
package main

import (
	"encoding/json"
	"flag"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sillypears/condor-standings/src/log"
	"github.com/sillypears/condor-standings/src/models"

	"regexp"
)

var (
	staticContent = flag.String("staticPath", "../public/swagger-ui", "Path to folder with Swagger UI")
	apiurl        = flag.String("api", "http://127.0.0.1", "The base path URI of the API service")
)

// ReturnedTables creates a struct for json output
type ReturnedTables struct {
	EventNames []string `json:"eventNames"`
}

func httpAPI(c *gin.Context) {
	// Local variables

	w := c.Writer

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(json.RawMessage(apiDescriptionsJson))
	if err != nil {
		log.Error("Couldn't generate JSON")
		return
	}

	w.Write(jsonData)

}

// @Title Events
// @Description Lists all events found by name
// @Accept plain
// @Produce json
// @Success 200 {object} ReturnedTables
// @Failure 404 {object} APIError "No Events Found"
// @Router /api/event [get]
func httpEventDocAPI(c *gin.Context) {
	// Local variables
	w := c.Writer

	returnedTables := make([]string, 0)

	foundTables, err := db.Tables.GetTables()
	if err != nil {
		log.Error("Could not get tables: ", err)
	}

	for _, tname := range foundTables {
		returnedTables = append(returnedTables, tname.TableName.String)
	}

	var TablesReturned ReturnedTables
	TablesReturned.EventNames = returnedTables

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.MarshalIndent(TablesReturned, "", "\t")
	if err != nil {
		log.Error("Couldn't generate JSON")
		//w.Write([]byte("Please search for a user"))
		return
	}
	w.Write(jsonData)

}

// @Title Event Listing
// @Description Lists everything found for the event
// @Accept plain
// @Produce json
// @Param event	path	string	true	"Event Name"
// @Success 200 {object} Event
// @Failure 404 {object} APIError "Event not found"
// @Router /api/event/{event} [get]
// httpEventAPI gets listings for the events
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
	if groups == true {
		foundEvent, _ = db.EventAPI.GetEventInfoGroups(event)
	} else {
		foundEvent, _ = db.EventAPI.GetEventInfo(event)
	}
	if err != nil {
		log.Error("Couldn't get event info: ", err)
	}
	season, err := regexp.MatchString("season_*", event)
	if err != nil {
		log.Error("Regex is bad")
	}

	editedEvent := models.Event{
		EventName: foundEvent.EventName,
	}
	var parts []models.Participant

	if season == true {
		for _, participant := range foundEvent.Participants {
			part := models.Participant{
				DiscordUsername: participant.DiscordUsername,
				TwitchUsername:  participant.TwitchUsername,
				EventPoints:     participant.EventPoints,
				EventPlayed:     participant.EventPlayed,
			}
			if event == "season_7" {
				part.TierName = s7tiers[strings.ToLower(participant.TwitchUsername)]
			}
			parts = append(parts, part)
		}
		editedEvent.Participants = parts
		foundEvent = editedEvent
	}

	jsonData, err := json.MarshalIndent(foundEvent, "", "\t")
	if err != nil {
		log.Error("Couldn't generate JSON")
		//w.Write([]byte("Please search for a user"))
		return
	}

	w.Write(jsonData)
}

// @Title Team Results Listing
// @Description Lists everything found for the season 7 teams
// @Accept plain
// @Produce json
// @Success 200 {object} models.Result
// @Failure 404 {object} APIError "Nothing found"
// @Router /api/teamresults [get]
// httpTeamAPI gets the team results and returns json
func httpTeamAPI(c *gin.Context) {
	w := c.Writer

	w.Header().Set("Content-Type", "application/json")

	results, err := db.TeamAPI.GetResults()
	if err != nil {
		log.Error("Couldn't get team info, ", err)
	}

	jsonData, err := json.MarshalIndent(results, "", "\t")
	if err != nil {
		log.Error("Couldn't parse results: ", err)
		return
	}
	w.Write(jsonData)
}
