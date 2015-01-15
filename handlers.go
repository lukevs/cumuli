package main 

import (
    "encoding/json"
    "html/template"
    // "log"
    "net/http"
    "os"
    "strings"
    "time"
)

var templates map[string]*template.Template

// MainHandler handles the route '/'.
func MainHandler(rw http.ResponseWriter, r *http.Request) {

    // Get the query parameter for u
    uParam := r.FormValue("u")
    if uParam == "" {

        // Render the splash page
        renderTemplate(rw, "splash.html", nil)
        return

    }

    var filename string = `./static/json/` + strings.Replace(uParam, " ", "+", -1) + `.json`

    // Handle file doesn't exist
    if _, err := os.Stat(filename); os.IsNotExist(err) {

        // Split fParam into individual users
        users := strings.Split(uParam, " ")

        // Get the shared followings among the users
        result := GetSharedFollowings(&users)

        // JSON marshal the result
        out, err := json.Marshal(*result)
        if err != nil {
            http.Error(rw, err.Error(), http.StatusInternalServerError)
            return
        }

        // Create a file <uParam>.json
        f, err := os.Create(filename)
        if err != nil {
            http.Error(rw, err.Error(), http.StatusInternalServerError)
            return
        }
        defer f.Close()

        // Store the JSON in it
        _, err = f.Write(out)
        if err != nil {
            http.Error(rw, err.Error(), http.StatusInternalServerError)
            return
        }

        // Hold on to the result for EXPIRE_TIME seconds
        // to cover user refreshes
        go func (filename string) {
            time.Sleep(time.Second * EXPIRE_TIME)
            os.Remove(filename)
        } (filename)
    }

    // Render the package
    renderTemplate(rw, "index.html", filename)
}

// StaticHandler handles the static assets of the app.
func StaticHandler(rw http.ResponseWriter, r *http.Request) {

    suffixes := []string{".json", ".css", ".otf", ".ico"}

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