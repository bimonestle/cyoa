package cyoa

import (
	"log"
	"strings"
	"encoding/json"
	"io"
	"net/http"
	"html/template"
)

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTemplate))
}

var tpl *template.Template

var defaultHandlerTemplate = `
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
		<li><a href="/{{.Chapter}}">{{.Text}}</a></li>
	{{end}}
    </ul>
</body>
</html>`

func NewHandler(s Story) http.Handler {
	return handler{s}
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)

	// If there is not specific path matched
	// It will just take you to the intro chapter
	if path == "" || path == "/" {
		path = "/intro"
	}
	// It says give us the path without the prefix "/"
	// "/intro" --> "intro"
	path = path[1:]

	// The first argument returned is gonna be the actual object stored in the map
	// the second argument is whether or not we found that actual key in the map
	if chapter, ok := h.s[path]; ok {
		err := tpl.Execute(w, chapter)
		if err != nil {
			// This log err part is only visible to the developer. Not the end user
			log.Printf("%v", err)
			
			// The following code is fine for development
			// So that we know what's happening with the error
			// http.Error(w, err, http.StatusFound)

			// No need to display the actual error to the end user
			// Because sometimes the err contains sensitive information
			// Such as password, etc. So in that case display some generic
			// error message
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter Not Found", http.StatusNotFound)
}

// Decode the opened/chosen json file and stores it into story variable
func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

type Story map[string]Chapter

type Chapter struct {
	Title     string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options   []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Chapter  string `json:"chapter"`
}
