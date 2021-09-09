package main

import (
	"chucknorris/internal/api"
	"flag"
	"fmt"
	"log"
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
		log.Fatal("Expected commands `random` and `dump`")
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
		log.Fatal("Expected commands `random` and `dump`")
	}
}

func checkError(err error) {
	// Check occurance of error
	if err != nil {
		log.Fatal(err)
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
			fmt.Println(err)
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

		retryNumber = 0
		jokes[jokeData.Id] = *jokeData
	}

	// Create file `category`.txt
	file, err := os.Create(category + ".txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Write jokes to file
	for _, jokeData := range jokes {
		file.WriteString(jokeData.Value + "\n")
	}
}
