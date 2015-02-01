// handlers.go contains the handler functions for cumuli

package main 

import (
    "html/template"
    "net/http"
    "path"
    "strings"
    "time"

    "github.com/garyburd/redigo/redis"
    "github.com/lkvnstrs/cumuli/networkmapper"
)

var templates map[string]*template.Template
const EXPIRE_TIME = 60 // in seconds

// MainHandler handles the route '/'.
func MainHandler(rw http.ResponseWriter, r *http.Request) {

    renderTemplate(rw, "splash.html", nil)
}

// UserHandler handles the display of D3 graphs for a given set of users
// at the route '/u/'.
func UserHandler(rw http.ResponseWriter, r *http.Request) {

    // Get the path base
    var jsonPath string = `/json/` + path.Base(r.URL.Path)
    jsonPath = strings.Trim(jsonPath, "+")

    // Render the page
    renderTemplate(rw, "index.html", jsonPath)
}

// AboutHandler handles the about page.
func AboutHandler(rw http.ResponseWriter, r *http.Request) {

    // Render the about page
    renderTemplate(rw, "about.html", nil)    
}

// JSONHandler handles the generation and display of JSON for D3 at the
// the route '/json/'.
func JSONHandler(rw http.ResponseWriter, r *http.Request) {

    var js []byte

    conn := pool.Get()
    defer conn.Close()

    // Get the path base
    key := path.Base(r.URL.Path)

    if key == "" {
        rw.Header().Set("Content-Type", "application/json")
        rw.Write([]byte{})
    }

    js, err := redis.Bytes(conn.Do("GET", key))

    // Handle key doesn't exist
    if err == redis.ErrNil {

        users := strings.Split(key, "+")

        js, err = networkmapper.BuildNetworkMap(n, users[0:])
        if err != nil {
            http.Error(rw, err.Error(), http.StatusInternalServerError)
            return
        }

        // Store the result
         _, err = conn.Do("SET", key, js)
        if err != nil {
            http.Error(rw, err.Error(), http.StatusInternalServerError)
            return
        }

        go func (key string) {
            time.Sleep(time.Second * EXPIRE_TIME)
            conn := pool.Get()
            defer conn.Close()
            conn.Do("DEL", key)
        } (key)
        

    } else if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
        return
    }

    // Render the JSON
    rw.Header().Set("Content-Type", "application/json")
    rw.Write(js)
}

// StaticHandler handles the static assets of the app.
func StaticHandler(rw http.ResponseWriter, r *http.Request) {

    suffixes := []string{".json", ".css", ".otf", ".ico", ".js"}

    for _, s := range suffixes {
        if strings.HasSuffix(r.URL.Path, s) {
            http.ServeFile(rw, r, r.URL.Path[1:])
            return
        }
    }

    http.NotFound(rw, r)
}

/* Helpers */

// renderTemplate is used to avoid code repetition for calling the 
func renderTemplate(rw http.ResponseWriter, filename string, data interface{}) {
    if err := templates[filename].ExecuteTemplate(rw, "base", data); err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}