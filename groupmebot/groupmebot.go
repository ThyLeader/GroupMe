// Package groupmebot handles the incoming and outgoing message functionality for gogroupmebot
package groupmebot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// The API base URL for GroupMe
const baseURL = "https://api.groupme.com/v3/bots/post"

// IncomingMessage is the data values for an incoming GroupMe message
type IncomingMessage struct {
	Id       string `json:"id"`
	User_id  string `json:"user_id"`
	Group_id string `json:"group_id"`
	Name     string `json:"name"`
	Text     string `json:"text"`
}

// OutgoingMessage is the data values for an outgoing GroupMe message
type OutgoingMessage struct {
	BotID string `json:"bot_id"`
	Text  string `json:"text"`
}

// PostMessage posts the given OutgoingMessage to the group the bot belongs to
func PostMessage(message OutgoingMessage) {
	// Parsing the Outgoingmessage and prepping a POST
	outgoing, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Marshaling of the OutgoingMessage failed:", err)
	}
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(outgoing))
	if err != nil {
		fmt.Println("Creating a new POST request failed:", err)
	}

	// Setting the headers required by the GroupMe API
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Sending the bot's reply
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Posting an outgoing message failed:", err)
	}
	resp.Body.Close()
}
