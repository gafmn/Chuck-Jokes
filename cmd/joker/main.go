package main

import (
	"chucknorris/internal/api"
	"flag"
	"fmt"
	"os"
	"sync"
)

func main() {
	// Entrypoint of program
	_ = flag.NewFlagSet("random", flag.ExitOnError)

	// Provide number of jokes to be dumped to files
	dumpCmd := flag.NewFlagSet("dump", flag.ExitOnError)
	dumpNumber := dumpCmd.Int("n", 5, "Number of jokes")

	// Check correct amount subcommands
	if len(os.Args) < 2 {
		fmt.Println("Expected commands `random` and `dump`")
		os.Exit(1)
	}

	// Check what command
	switch os.Args[1] {

	case "random":
		joke, err := api.GetRandomJoke()
		checkError(err)

		fmt.Println(joke)

	case "dump":
		dumpCmd.Parse(os.Args[2:])

		// Get all categories
		categories, err := api.GetCategoryList()
		checkError(err)

		// Save n jokes of categories
		saveNJokes(categories, *dumpNumber)

	default:
		// Case in incorrect commands
		fmt.Println("Expected commands `random` and `dump`")
		os.Exit(1)
	}
}

func checkError(err error) {
	// Check occurance of error
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func saveNJokes(categories []string, n int) {
	// Save N jokes from all categories
	var wg sync.WaitGroup

	// Iterate over categories and fetch random jokes
	for _, category := range categories {
		wg.Add(1)
		go saveCategory(category, n, &wg)
	}

	wg.Wait()
}

func saveCategory(category string, n int, wg *sync.WaitGroup) {
	// Save category jokes to file
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
