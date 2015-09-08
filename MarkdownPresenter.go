package main
import (
  "log"
  "fmt"
  "net/http"
  "io/ioutil"
  "github.com/russross/blackfriday"
  "github.com/microcosm-cc/bluemonday"
)

type Page struct {
    Title string
    Body  []byte
}

var port = "8080"
var hostname = "http://localhost:" + port
var folder = "./"

func loadPage(title string) (*Page, error) {
    filename := title + ".md"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/view/"):]
    p, _ := loadPage(title)
    content := blackfriday.MarkdownCommon(p.Body)
    html := bluemonday.UGCPolicy().SanitizeBytes(content)
    fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, html)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
  log.Printf("indexHandler")

  files, _ := ioutil.ReadDir(folder)

  fmt.Fprintf(w, "<h1>Folder list:</h>\n<ul>\n")
  for _, file := range files {
    if file.IsDir() {
      log.Printf("Folder :" + file.Name())
      fmt.Fprintf(w, "<li><a href=\"" + hostname + "/view/" + file.Name() + "/README\">" + file.Name() + "</a></li>\n")
    }
  }
  fmt.Fprintf(w, "</ul>\n")
}

func main() {
    log.Printf("About to listen on port: " + port)
    http.HandleFunc("/view/", viewHandler)
    http.HandleFunc("/index/", indexHandler)
    http.ListenAndServe(":" + port, nil)
}
