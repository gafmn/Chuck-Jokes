package main

import (
	"chucknorris/internal/api"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

// Entrypoint of program
func main() {
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

		categories, err := api.GetCategoryList()
		checkError(err)

		saveNJokes(categories, *dumpNumber)

	default:
		log.Fatal("Expected commands `random` and `dump`")
	}
}

// Check occurance of error
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Save N jokes from all categories
func saveNJokes(categories []string, n int) {
	var wg sync.WaitGroup

	for _, category := range categories {
		wg.Add(1)
		go saveCategory(category, n, &wg)
	}

	wg.Wait()
}

// Save category jokes to file
func saveCategory(category string, n int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Declare map of jokes and number of tries to get random joke of category
	jokes := make(map[string]api.JokeData)
	retryNumber := 0

	for len(jokes) < n {
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
