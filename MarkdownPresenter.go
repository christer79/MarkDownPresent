package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

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
	log.Print(body)
	lines := strings.Split(string(body), "\n")
	author := ""
	created := ""
	defaultBackground := ""
	if mdFormat.IsCommentedLine(lines[0]) {
		log.Println("Found comment")
		author = mdFormat.ExtractCommentDataFiled(lines[0], "Author")
		created = mdFormat.ExtractCommentDataFiled(lines[0], "Created")
		defaultBackground = mdFormat.ExtractCommentDataFiled(lines[0], "Background")
		lines = lines[1:]
	}
	log.Println("Default background: " + defaultBackground)
	log.Println("Author: " + author)
	log.Println("Created date: " + created)
	var page []byte
	var pages []Slide
	for _, line := range lines {
		log.Println(line)
		if mdFormat.IsCommentedLine(line) {
			timeout, _ := strconv.Atoi(mdFormat.ExtractCommentDataFiled(line, "Timeout"))
			comment := mdFormat.ExtractCommentDataFiled(line, "Comment")
			background := mdFormat.ExtractCommentDataFiled(line, "Background")
			if background == "" {
				background = defaultBackground
			}
			pages = append(pages, Slide{Timeout: timeout, Comment: comment, Background: background, Body: page, NextSlideNr: len(pages) + 1})
			page = []byte{}
		} else {
			page = append(page, []byte(line+"\n")...)
		}
	}
	//TODO: Reload only if syntax is OK!
	pages[len(pages)-1].NextSlideNr = 0
	presentation := Presentation{Filename: filename, Pages: pages}

	return presentation, err
}

func loadPage(filename string, slide int) (*Slide, error) {
	presentation, _ := loadPresentation(filename)
	return &presentation.Pages[slide], nil
}

func getFileName(URL string) string {
	re := regexp.MustCompile("/view/(.*)/[0-9]*$")
	matches := re.FindStringSubmatch(URL)
	filename := matches[1]
	log.Println("Regexp match filename #: " + filename)
	return filename
}

func getSlideNr(URL string) int {
	log.Println("URL: \"" + URL + "\"")
	re := regexp.MustCompile("/([0-9]*)$")
	matches := re.FindStringSubmatch(URL)
	slide, err := strconv.Atoi(matches[1])
	if err != nil {
		panic(err)
	}
	log.Println("Regexp match slide #: " + strconv.Itoa(slide))
	return slide
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	slideNr := getSlideNr(r.URL.Path)
	log.Println("Slidenr: " + strconv.Itoa(slideNr))
	filename := getFileName(r.URL.Path)
	log.Println("filename: " + filename)

	p, _ := loadPage(filename, slideNr)
	html := mdFormat.MarkDownToHTML(p.Body)

	var timeoutString = ""
	if p.Timeout > 0 {
		timeoutString = "<script>setTimeout(function(){ location.href = \"" + strconv.Itoa(p.NextSlideNr) + "\"; }, " + strconv.Itoa(p.Timeout*1000) + ");</script>"
	}
	ajaxString := "<script src=\"https://ajax.googleapis.com/ajax/libs/jquery/1.11.3/jquery.min.js\"></script>"
	mouseClickString := "<script>$(document).ready(function(){$(\"body\").click(function(){location.href = \"" + strconv.Itoa(p.NextSlideNr) + "\";});});</script>"
	//cssString := "<style>h1 {color:red;}</style>"
	//cssString := "" //"<link rel=\"stylesheet\" href=\"style.css\">"
	body, _ := ioutil.ReadFile("style.css")
	cssString := "<style>" + string(body) + "</style>"

	fmt.Fprintf(w, cssString+ajaxString+timeoutString+mouseClickString+"<body background=\""+p.Background+"\"><div>%s</div>", html)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("indexHandler")

	//TODO: Select only files
	//TODO: Recurse into subfolders
	//files, _ := ioutil.ReadDir(folder)
	files := fileTools.FindFiles(folder, "*.md")
	fmt.Fprintf(w, "<h1>Folder list:</h>\n<ul>\n")
	for _, file := range files {
		log.Printf("Folder :" + file)
		// TODO: Print number of pagaes and links to each page on index page
		fmt.Fprintf(w, "<li><a href=\""+hostname+"/view/"+file+"/0\">"+file+"</a></li>\n")
	}
	fmt.Fprintf(w, "</ul>\n")
}

func main() {
	log.Printf("About to listen on port: " + port)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/index/", indexHandler)
	http.ListenAndServe(":"+port, nil)
}
