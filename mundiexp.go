package mundiexp

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"io/ioutil"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
)

const (
	imgDir      = "public/img/"
	markdownDir = "markdown/"
)

var tmpl *template.Template

type templateData struct {
	ID            int
	ImageFilePath string
	HTML          template.HTML
}

func home(w http.ResponseWriter, req *http.Request) {
	td := make([]templateData, 0)
	files, err := ioutil.ReadDir(imgDir)
	if err != nil {
		log.Fatalf("couldn't get file list from public/img/: %v", err)
	}
	for i, file := range files {
		imageFileName := file.Name()
		ext := filepath.Ext(imageFileName)
		// check for extension
		if ext != ".jpg" && ext != ".png" {
			log.Fatalf("unexpected file extension %v for file %v in public/img/. jpg and png are the only acceptable ones.", ext, imageFileName)
		}

		td = append(td, templateData{ID: i, ImageFilePath: imgDir + imageFileName})

		mdFileName := imageFileName[:len(imageFileName)-len(filepath.Ext(imageFileName))] + ".md"
		mdFilePath := markdownDir + mdFileName
		md, err := ioutil.ReadFile(mdFilePath)
		if err != nil {
			log.Printf("couldn't find corresponding markdown file for image %v. continuing", td[i].ImageFilePath)
			continue
		}
		html := blackfriday.MarkdownCommon(md)
		td[i].HTML = template.HTML(html)
	}

	err = tmpl.Execute(w, td)
	if err != nil {
		log.Fatalf("couldn't execute template: %v", err)
	}
}

func init() {
	// parse template

	var err error
	tmpl, err = template.ParseFiles("index.html")
	if err != nil {
		log.Fatalf("couldn't parse template: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", home).Methods("GET")
	http.Handle("/", r)
}
