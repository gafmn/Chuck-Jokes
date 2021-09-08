package main

import (
	"chucknorris/internal/api"
	"flag"
	"fmt"
	"os"
	"sync"
)

func main() {
	fmt.Println("Hello there!")
	_ = flag.NewFlagSet("random", flag.ExitOnError)

	dumpCmd := flag.NewFlagSet("dump", flag.ExitOnError)
	dumpNumber := dumpCmd.Int("n", 5, "Number of jokes")

	if len(os.Args) < 2 {
		fmt.Println("Expected commands `random` and `dump`")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "random":
		joke, err := api.GetRandomJoke()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println(joke)

	case "dump":
		dumpCmd.Parse(os.Args[2:])

		categories, err := api.GetCategoryList()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		saveNJokes(categories, *dumpNumber)

	default:
		fmt.Println("Expected commands `random` and `dump`")
		os.Exit(1)
	}
}

func saveNJokes(categories []string, n int) error {
	var wg sync.WaitGroup

	// Iterate over categories and fetch random jokes
	for _, category := range categories {
		wg.Add(1)
		go saveCategory(category, n, &wg)
	}

	wg.Wait()

	return nil
}

func saveCategory(category string, n int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Declare map of jokes and number of tries to get random joke of category
	jokes := make(map[string]api.JokeData)
	retryNumber := 0

	for len(jokes) < n {
		// Fetch random joke of category
		jokeData, err := api.GetCategoryRandomJoke(category)

		if err != nil {
			fmt.Println(err.Error())
		}

		// Check uniqueness
		if _, ok := jokes[jokeData.Id]; ok {
			if retryNumber < len(jokes)*5 {
				retryNumber += 1
				continue
			}
			fmt.Printf("WARNING: In category %v probably jokes less than %d\n", category, n)
			break
		}

		// Update data
		retryNumber = 0
		jokes[jokeData.Id] = *jokeData
	}

	// Create file `category`.txt
	file, err := os.Create(category + ".txt")
	defer file.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	// Write jokes to file
	for _, jokeData := range jokes {
		file.WriteString(jokeData.Value + "\n")
	}
}
