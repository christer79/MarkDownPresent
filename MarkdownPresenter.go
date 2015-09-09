package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

type Page struct {
	Timeout     int
	Comment     string
	Background  string
	Description string
	NextSlideNr int
	Body        []byte
}

type Presentation struct {
	Filename   string
	Author     string
	Background string
	Pages      []Page
}

var port = "8000"
var hostname = "http://localhost:" + port
var folder = "./"

func isCommentedLine(line string) bool {
	re := regexp.MustCompile("[[//]]: # ")
	return re.MatchString(line)
}

func extractCommentDataFiled(line string, label string) string {
	commentRe := regexp.MustCompile(label + ": \"([^\"]*)\"")
	value := commentRe.FindAllStringSubmatch(line, -1)[0][1]
	log.Println("Found " + label + ": \"" + value + "\"")
	return value
}

func loadPresentation(filename string) (Presentation, error) {
	//TODO: Add one comment line with global settings
	//TODO: use sha1sum to know if file changed
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(body), "\n")
	author := ""
	created := ""
	defaultBackground := ""
	if isCommentedLine(lines[0]) {
		author = extractCommentDataFiled(lines[0], "Author")
		created = extractCommentDataFiled(lines[0], "Created")
		defaultBackground = extractCommentDataFiled(lines[0], "Background")
		lines = lines[1:]
	}
	log.Println("Default background: " + defaultBackground)
	log.Println("Author: " + author)
	log.Println("Created date: " + created)
	var page []byte
	var pages []Page
	for _, line := range lines {
		if isCommentedLine(line) {
			timeout, _ := strconv.Atoi(extractCommentDataFiled(line, "Timeout"))
			comment := extractCommentDataFiled(line, "Comment")
			background := extractCommentDataFiled(line, "Background")

			pages = append(pages, Page{Timeout: timeout, Comment: comment, Background: background, Body: page, NextSlideNr: len(pages) + 1})
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

func loadPage(filename string, slide int) (*Page, error) {
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
	content := blackfriday.MarkdownCommon(p.Body)
	html := bluemonday.UGCPolicy().SanitizeBytes(content)

	var timeoutString = ""
	if p.Timeout > 0 {
		timeoutString = "<script>setTimeout(function(){ location.href = \"" + strconv.Itoa(p.NextSlideNr) + "\"; }, " + strconv.Itoa(p.Timeout*1000) + ");</script>"
	}
	fmt.Fprintf(w, "<script src=\"https://ajax.googleapis.com/ajax/libs/jquery/1.11.3/jquery.min.js\"></script>"+timeoutString+"<script>$(document).ready(function(){$(\"body\").click(function(){location.href = \""+strconv.Itoa(p.NextSlideNr)+"\";});});</script><body background=\""+p.Background+"\"><h1>%s</h1><div>%s</div>", p.Comment, html)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("indexHandler")

	//TODO: Select only files
	//TODO: Recurse into subfolders
	files, _ := ioutil.ReadDir(folder)

	fmt.Fprintf(w, "<h1>Folder list:</h>\n<ul>\n")
	for _, file := range files {
		if file.IsDir() {
			log.Printf("Folder :" + file.Name())
			// TODO: Print number of pagaes and links to each page on index page
			fmt.Fprintf(w, "<li><a href=\""+hostname+"/view/"+file.Name()+"/README.md/0\">"+file.Name()+"</a></li>\n")
		}
	}
	fmt.Fprintf(w, "</ul>\n")
}

func main() {
	log.Printf("About to listen on port: " + port)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/index/", indexHandler)
	http.ListenAndServe(":"+port, nil)
}
