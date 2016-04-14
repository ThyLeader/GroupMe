// Package reddit returns posts via the reddit API
package reddit

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
