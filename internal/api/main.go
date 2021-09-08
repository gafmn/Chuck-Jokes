package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type JokeData struct {
	IconUrl string `json:"icon_url"`
	Id      string
	Url     string
	Value   string
}

func GetRandomJoke() (string, error) {
	// Fetch random joke from API
	response, err := http.Get("https://api.chucknorris.io/jokes/random")
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to get random joke")
	}

	// Get body of response on success
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	// Parse body
	var jokeData JokeData
	json.Unmarshal([]byte(body), &jokeData)

	return jokeData.Value, nil
}

func GetCategoryList() ([]string, error) {
	// Fetch all categories
	response, err := http.Get("https://api.chucknorris.io/jokes/categories")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to get list of categories")
	}

	// Get body of response on success
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// Parse body
	var categories []string
	json.Unmarshal([]byte(body), &categories)

	return categories, nil
}

func GetCategoryRandomJoke(category string) (*JokeData, error) {
	// Fetch random joke by category

	// Geneerate query param for joke request
	request, _ := http.NewRequest("GET", "https://api.chucknorris.io/jokes/random", nil)
	query := request.URL.Query()
	query.Add("category", category)
	request.URL.RawQuery = query.Encode()

	response, err := http.Get(request.URL.String())
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed get joke of category %v", category)
	}

	// Get body of response on success
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// Parse body
	var jokeData JokeData
	json.Unmarshal([]byte(body), &jokeData)

	return &jokeData, nil
}
