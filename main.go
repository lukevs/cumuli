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
//  - Add a search bar
//      - http://jsbin.com/iyewas/73/edit?html,js,output
//      - http://blogs.msdn.com/b/murads/archive/2013/02/20/using-jquery-ui-autocomplete-with-the-rest-api-to-get-search-results-in-the-search-box.aspx
//      - http://stackoverflow.com/questions/14083272/how-to-make-a-tags-box-using-jquery-with-text-input-field-tags-separated-by
//  - Add a loading gif

package main 

import (
    "html/template"
    "io/ioutil"
    "log"
    "net/http"
    "os"
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
    port := GetPort()
    log.Println("Running on port ", port)
    http.ListenAndServe(port, nil)
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

// for Heroku
func GetPort() string {
        var port = os.Getenv("PORT")
        // Set a default port if there is nothing in the environment
        if port == "" {
                port = "8080"
                log.Println("INFO: No PORT environment variable detected, defaulting to " + port)
        }
        return ":" + port
}