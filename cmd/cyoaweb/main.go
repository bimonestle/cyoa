package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bimonestle/go-exercise-projects/03.Choose-your-own-adventure/cyoa"
)

func main() {
	port := flag.Int("port", 3000, "The port to start the CYOA web application on")
	filename := flag.String("file", "gopher.json", "The json file with the CYOA story")
	flag.Parse()
	fmt.Printf("Using the story in %s.\n", *filename)

	// Get the file (the json file) from the localdisk
	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}

	story, err := cyoa.JsonStory(f)
	if err != nil {
		panic(err)
	}

	h := cyoa.NewHandler(story)
	fmt.Printf("Starting the sever on port at %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
	// fmt.Printf("%+v\n", h)
}