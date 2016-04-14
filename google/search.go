// Package search handles querying Google's Custom Search API
package search

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
