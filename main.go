// cumuli - A followings visualizer for SoundCloud

// To do:
//  - Add http-based error handling universally
//  - Add dynamic circle sizing
//  - Add dynamic resizing w/o refresh
//  - Add logging
//  - Add the About page
//  - Link the JS
//  - Break followings.go into more sub-functions
//  - Serve a favicon at /favicon.ico
//  - Patch MainHandler getting called 2/3 times

package main 

import (
    "html/template"
    "io/ioutil"
    "log"
    "net/http"
    "path"
)

const TEMPLATES_DIR = `./templates`

func init() {

    log.SetFlags(log.LstdFlags | log.Lmicroseconds)

    loadTemplates()

    // Routes
    http.HandleFunc("/", MainHandler)
    http.HandleFunc("/static/", StaticHandler)    
}

func main() {
    log.Println("Running on port :8080")
    http.ListenAndServe(":8080", nil)
}

func loadTemplates() {
    if templates == nil {
        templates = make(map[string]*template.Template)
    }

    // Import each file as an extension of base.html.
    files, _ := ioutil.ReadDir(TEMPLATES_DIR)
    base := path.Join(TEMPLATES_DIR, "base.html")

    for _, f := range files {
        if f.Name() != "base.html" {
            templates[f.Name()] = template.Must(template.ParseFiles(path.Join(TEMPLATES_DIR, f.Name()), base))
        }
    }
}