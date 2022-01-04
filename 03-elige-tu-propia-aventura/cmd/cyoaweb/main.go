package main

import (
	"cyoa"
	"flag"
	"fmt"
	"os"
)

const ERR_NOT_ABLE_TO_OPEN_FILE = 1

func main() {
	filename := flag.String("file", "gopher.json", "The JSON file with the story")
	flag.Parse()
	fmt.Printf("Using the story from file %s.\n", *filename)

	f, err := os.Open(*filename)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(ERR_NOT_ABLE_TO_OPEN_FILE)
	}

	story, err := cyoa.JsonStory(f)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", story)
}
