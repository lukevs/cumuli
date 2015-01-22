// cumuli - A followings visualizer for SoundCloud

// To do:
//  - Add http-based error handling universally & logging
//  - Add dynamic circle sizing
//  - Add dynamic resizing w/o refresh
//  - Add the About page
//  - Break followings.go into more sub-functions
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

var (
    clientId string
    redisPort = ":6379"
)

func init() {

    // Set log flags
    log.SetFlags(log.LstdFlags | log.Lmicroseconds)

    // Load templates
    loadTemplates()

    clientId = GetClientId()

    // Routes
    http.HandleFunc("/", MainHandler)
    http.HandleFunc("/u/", UserHandler)
    http.HandleFunc("/json/", JSONHandler)
    http.HandleFunc("/static/", StaticHandler)    
}

func main() {

    // Get port
    port := GetPort()

    log.Println("Running on port ", port)
    http.ListenAndServe(port, nil)
}

// loadTemplates loads all of the templates in TEMPLATES_DIR to be served.
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

// GetPort gets a PORT env if set and returns 8080 otherwise.
func GetPort() string {
        var port = os.Getenv("PORT")
        // Set a default port if there is nothing in the environment
        if port == "" {
                port = "8080"
                log.Println("INFO: No PORT environment variable detected, defaulting to " + port)
        }
        return ":" + port
}

// GetClientId gets the Soundcloud API client id.
func GetClientId() string {
    ci := os.Getenv("SC_CLIENT_ID")
    if ci == "" {
        log.Fatal("You forgot SC_CLIENT_ID")
    }
    return ci
}
