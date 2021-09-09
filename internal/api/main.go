package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type JokeData struct {
	IconUrl string `json:"icon_url"`
	Id      string
	Url     string
	Value   string
}

var URL = "https://api.chucknorris.io/jokes"

// Fetch random joke from API
func GetRandomJoke() (string, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	response, err := client.Get(URL + "/random")
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to get random joke")
	}

	// Get body of response on success
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	var jokeData JokeData
	if err := json.Unmarshal([]byte(body), &jokeData); err != nil {
		fmt.Println(err)
		return "", err
	}

	return jokeData.Value, nil
}

// Fetch all categories
func GetCategoryList() ([]string, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	response, err := client.Get(URL + "/categories")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to get list of categories")
	}

	// Get body of response on success
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var categories []string
	if err := json.Unmarshal([]byte(body), &categories); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return categories, nil
}

// Fetch random joke by category
func GetCategoryRandomJoke(category string) (*JokeData, error) {

	// Geneerate query param for joke request
	request, _ := http.NewRequest("GET", URL+"/random", nil)
	query := request.URL.Query()
	query.Add("category", category)
	request.URL.RawQuery = query.Encode()

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	response, err := client.Get(request.URL.String())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed get joke of category %v", category)
	}

	// Get body of response on success
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var jokeData JokeData
	if err := json.Unmarshal([]byte(body), &jokeData); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &jokeData, nil
}
