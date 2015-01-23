// handlers.go contains the handler functions for cumuli

package main 

import (
    "encoding/json"
    "html/template"
    //"log"
    "net/http"
    "path"
    "strings"
    "time"

    "github.com/garyburd/redigo/redis"
)

var templates map[string]*template.Template

// MainHandler handles the route '/'.
func MainHandler(rw http.ResponseWriter, r *http.Request) {

    renderTemplate(rw, "splash.html", nil)
}

// UserHandler handles the display of D3 graphs for a given set of users
// at the route '/u/'.
func UserHandler(rw http.ResponseWriter, r *http.Request) {

    // Get the path base
    var jsonPath string = `/json/` + path.Base(r.URL.Path)

    // Render the page
    renderTemplate(rw, "index.html", jsonPath)
}

// JSONHandler handles the generation and display of JSON for D3 at the
// the route '/json/'.
func JSONHandler(rw http.ResponseWriter, r *http.Request) {

    var js []byte

    // Create a connection to Redis
    c := pool.Get()
    defer c.Close()

    // Get the path base
    key := path.Base(r.URL.Path)

    if key == "" {
        rw.Header().Set("Content-Type", "application/json")
        rw.Write([]byte{})
    }

    // Check the db for base as a key
    js, err := redis.Bytes(c.Do("GET", key))
    if err == redis.ErrNil {

        // Split fParam into individual users
        users := strings.Split(key, "+")

        // Get the shared followings among the users
        result := GetSharedFollowings(&users)    

        // JSON marshal the result
        js, err = json.Marshal(*result)
        if err != nil {
            http.Error(rw, err.Error(), http.StatusInternalServerError)
            return
        }

        // Store the JSON in Redis
        _, err = c.Do("SET", key, js)
        if err != nil {
            http.Error(rw, err.Error(), http.StatusInternalServerError)
            return
        }

        // Hold on to the result for EXPIRE_TIME seconds
        // to cover user refreshes
        go func (key string) {
            time.Sleep(time.Second * EXPIRE_TIME)
            c := pool.Get()
            defer c.Close()
            c.Do("DEL", key)
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