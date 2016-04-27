package main

import (
    "encoding/json"
    "fmt"
    "math/rand"
    "net/http"
    "os"
    "strings"
    "time"
    "bytes"
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
    var im IncomingMessage
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
        go BotAww(rand, Settings.RedditUserAgent, ch)
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
        om := OutgoingMessage{Settings.BotID, content}
        go PostMessage(om)
    }
}

// A redditResponse is a set of nested structs that translates the JSON response from reddit
// into a easier to work with data structure
type redditResponse struct {
	Data struct {
		Children []struct {
			Data redditPost
		}
	}
}

// A redditPost is the Title and URL of a post from reddit.
type redditPost struct {
	Title string
	URL   string
}

// BotAww queries /r/aww and returns a random selection from the top 100 posts of the last month
func BotAww(rand int, userAgent string, ch chan string) {
	subreddit := "https://www.reddit.com/r/aww/top.json?t=month&limit=100"
	client := &http.Client{}
	req, err := http.NewRequest("GET", subreddit, nil)
	if err != nil {
		fmt.Println("Starting a new request failed:", err)
	}

	// Adds the configured userAgent to the request string
	// This is necessary so that Reddit does not take action against the bot
	req.Header.Set("User-Agent", userAgent)

	// Do the request, and parse the JSON returned
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("The Reddit request failed:", err)
	}
	defer response.Body.Close()
	var r = new(redditResponse)
	err = json.NewDecoder(response.Body).Decode(&r)
	if err != nil {
		fmt.Println("JSON Decoding failed:", err)
	}

	// Creating an array of posts to pick at random
	posts := make([]redditPost, len(r.Data.Children))
	for i, child := range r.Data.Children {
		posts[i] = child.Data
	}
	var post = posts[rand]
	ch <- post.URL
}

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

/*
type googleResponse struct {
	Results []GoogleResult `json:"items"`
}

type GoogleResult struct {
	URL string `json:"link"`
}

func GoogleImage(query string, rand int, googleAPIKey string, googleCx string, ch chan string) {
	// Prepare the Google Search API request.
	client := &http.Client{}
	url := "https://www.googleapis.com/customsearch/v1?searchType=image&safe=off"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Starting a new request failed:", err)
	}

	// Adding the queries, key, and cx values to the request
	q := req.URL.Query()
	q.Set("q", query)
	q.Set("key", googleAPIKey)
	q.Set("cx", googleCx)
	req.URL.RawQuery = q.Encode()

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("The Google API request failed:", err)
	}
	defer response.Body.Close()

	// Creating the googleResponse data
	var data = new(googleResponse)
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		fmt.Println("Setting the googleResponse data failed:", err)
	}

	// Creating an array of responses to pick at random
	results := make([]GoogleResult, len(data.Results))
	for i, res := range data.Results {
		results[i] = res
	}
	result := results[rand]
	ch <- result.URL
}
*/