package main

import (
	"cyoa"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const ERR_NOT_ABLE_TO_OPEN_FILE = 1

func main() {
	port := flag.Int("port", 3000, "Port where the CYOA server listens")
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

	h := cyoa.NewHandler(story)
	fmt.Printf("Starting CYOA server on port %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
