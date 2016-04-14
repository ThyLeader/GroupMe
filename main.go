package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	//"github.com/ThyLeader/GroupMe/google"
	"github.com/ThyLeader/GroupMe/groupmebot"
	"github.com/ThyLeader/GroupMe/reddit"
)

var Settings struct {
	//GoogleAPIKey    string `json:"GoogleAPIKey"`
	//GoogleCx        string `json:"GoogleCX"`
	//BotUserID       string `json:"BotUserID"`
	BotID           string `json:"BotID"`
	RedditUserAgent string `json:"RedditUserAgent"`
	Port            string `json:"AppPort"`
}

// random returns a "random" integer - used for selecting a random result for a bot.
func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

// handler takes the incoming message and runs the correct reddit/groupmebot function
// depending on the text inputted by the user.
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New incoming message...")
	body := json.NewDecoder(r.Body)
	var im groupmebot.IncomingMessage
	body.Decode(&im)

	// Ignoring messages from the bot
	//if im.User_id == Settings.BotUserID {
	//	return
	//}

	// Converting the incomingMessage text to lowercase and running string comparisons to
	// get the requested contest
	var content string
	lText := strings.ToLower(im.Text)

	// If text contains !aww - run randomizer BotAww
	if strings.Contains(lText, "!aww") {
		ch := make(chan string)
		rand := random(0, 100)
		go reddit.BotAww(rand, Settings.RedditUserAgent, ch)
		content = <- ch
	}
	// If text contains !search - trim !search from the text, then pass it and settings to the GoogleImage func
	/*if strings.Contains(lText, "!search") {
		ch := make(chan string)
		searchText := strings.TrimPrefix(lText, "!search ")
		rand := random(0, 10)
		go search.GoogleImage(searchText, rand, Settings.GoogleAPIKey, Settings.GoogleCx, ch)
		content = <- ch
	}*/

	// If content is found, return it via bot message
	if len(content) > 0 {
		om := groupmebot.OutgoingMessage{Settings.BotID, content}
		go groupmebot.PostMessage(om)
	}
}

// Listening for traffic on port 8080 and passing to handler
func main() {
	// Setting configuration
	configFile, err := os.Open("conf.json")
	if err != nil {
		fmt.Println("Could not open config file!")
	}
	if err := json.NewDecoder(configFile).Decode(&Settings); err != nil {
		fmt.Println("Decoding settings failed:", err)
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe(Settings.Port, nil)
}
