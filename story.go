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
	// Parsing the path section

	// Get the path from the url and check to see
	// if the path equals to empty string or slash,
	// redirect it to the intro page
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	// get the path by slicing the /
	// "/intro" --> "intro"
	path = path[1:]

	// The first argument returned is going to be the actual object stored in the map; The chapter
	// The second argument is whether or not it actually find the key inside the map
	if chapter, ok := h.s[path]; ok {
		err := tpl.Execute(w, chapter)
		if err != nil {
			// Log the error to get to know what the actual error is about
			log.Printf("%v", err)

			// Sometimes the developer would print the error message as is
			// just to know what the actual error is and it's fine for
			// the development phase; not for build phase. re: the following code
			// http.Error(w, err, http.StatusInternalServerError)

			// The reason why the error message is quite generic because
			// sometimes it contains sensitive data within the returned error
			// for example is like password etc etc
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
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
	Chapter  string `json:"arc"`
}
