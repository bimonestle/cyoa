package main

import (
	"html/template"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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

	tpl:=template.Must(template.New("").Parse(storyTemplate))
	
	// This code below is for the default template
	// h := cyoa.NewHandler(story)

	// This code below is for the custom template
	// tpl := template.Must(template.New("").Parse("Hello!"))
	h := cyoa.NewHandler(story,
		cyoa.WithTemplate(tpl),
		cyoa.WithPathFunc(pathFn))

	mux:=http.NewServeMux()
	mux.Handle("/story/", h)

	fmt.Printf("Starting the sever on port at %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
	// fmt.Printf("%+v\n", h)
}

func pathFn(r *http.Request) string  {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}

var storyTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Choose Your Own Adventure</title>
</head>
<body>
    <h1>{{.Title}}</h1>
    {{range .Paragraphs}}
	<p>{{.}}</p>
    {{end}}
    <ul>
	{{range .Options}}
		<li><a href="/story/{{.Chapter}}">{{.Text}}</a></li>
	{{end}}
    </ul>
</body>
</html>`