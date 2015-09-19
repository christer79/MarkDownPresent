package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/christer79/MarkDownPresent/fileTools"
	"github.com/christer79/MarkDownPresent/mdFormat"
)

// Slide represents one slinge slide with infomration about time to show and content
type Slide struct {
	Timeout     int
	Comment     string
	Background  string
	Description string
	NextSlideNr int
	Style       string
	MarkDown    []byte
	Body        []byte
}

//Presentation represents a collection of Slides as well as information about author and creation date
type Presentation struct {
	Filename   string
	Author     string
	Background string
	Pages      []Slide
}

var port = "8000"
var hostname = "http://localhost:" + port
var folder = "./examples"

func loadPresentation(filename string) (Presentation, error) {
	//TODO: use sha1sum to know if file changed
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(body), "\n")
	author := ""
	created := ""
	defaultStyle := ""
	defaultBackground := ""
	if mdFormat.IsCommentedLine(lines[0]) {
		log.Println("Found comment")
		author = mdFormat.ExtractCommentDataFiled(lines[0], "Author", "")
		created = mdFormat.ExtractCommentDataFiled(lines[0], "Created", "")
		defaultStyle = mdFormat.ExtractCommentDataFiled(lines[0], "Style", "")
		defaultBackground = mdFormat.ExtractCommentDataFiled(lines[0], "Background", "")
		lines = lines[1:]
	}
	log.Println("Author: " + author)
	log.Println("Created date: " + created)

	var page []byte
	var slides []Slide
	for _, line := range lines {
		if mdFormat.IsCommentedLine(line) {
			timeout, _ := strconv.Atoi(mdFormat.ExtractCommentDataFiled(line, "Timeout", "20000000"))
			timeout = timeout * 1000
			if timeout == 0 {
				timeout = 300000000
			}
			comment := mdFormat.ExtractCommentDataFiled(line, "Comment", "")
			style := mdFormat.ExtractCommentDataFiled(line, "Style", defaultStyle)
			background := mdFormat.ExtractCommentDataFiled(line, "Background", defaultBackground)
			if background == "" {
				background = defaultBackground
			}
			body := mdFormat.MarkDownToHTML(page)
			slides = append(slides, Slide{Timeout: timeout, Style: style, Comment: comment, Background: background, Body: body, MarkDown: page, NextSlideNr: len(slides) + 1})
			page = []byte{}
		} else {
			page = append(page, []byte(line+"\n")...)
		}
	}
	//TODO: Reload only if syntax is OK!
	slides[len(slides)-1].NextSlideNr = 0
	presentation := Presentation{Filename: filename, Pages: slides}

	return presentation, err
}

func loadSlide(filename string, slide int) (*Slide, error) {
	presentation, _ := loadPresentation(filename)
	return &presentation.Pages[slide], nil
}

func getFileName(URL string) string {
	re := regexp.MustCompile("/view/(.*)/[0-9]*$")
	matches := re.FindStringSubmatch(URL)
	filename := matches[1]
	return filename
}

func getSlideNr(URL string) int {
	re := regexp.MustCompile("/([0-9]*)$")
	matches := re.FindStringSubmatch(URL)
	slide, err := strconv.Atoi(matches[1])
	if err != nil {
		panic(err)
	}
	return slide
}

func styleHandler(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/")
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(body))
}

func backgroundHandler(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/")
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(body))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	slideNr := getSlideNr(r.URL.Path)
	filename := getFileName(r.URL.Path)

	p, _ := loadSlide(filename, slideNr)
	t, err := template.ParseFiles("html.tmpl/view.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, p)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	files := fileTools.FindFiles(folder, "*.md")
	fmt.Fprintf(w, "<h1>Folder list:</h>\n<ul>\n")
	for _, file := range files {
		// TODO: Print number of pagaes and links to each page on index page
		fmt.Fprintf(w, "<li><a href=\""+hostname+"/view/"+file+"/0\">"+file+"</a></li>\n")
	}
	fmt.Fprintf(w, "</ul>\n")
}

func main() {
	log.Printf("About to listen on port: " + port)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/index/", indexHandler)
	http.HandleFunc("/style/", styleHandler)
	http.Handle("/background/", http.StripPrefix("/background/", http.FileServer(http.Dir("background"))))
	http.ListenAndServe(":"+port, nil)
}
