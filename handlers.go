package main 

import (
    "encoding/json"
    "fmt"
    "html/template"
    "io/ioutil"
    "net/http"
    "os"
    "path"
    "strings"
    "time"
)

// MainHandler handles the route '/'.
func MainHandler(rw http.ResponseWriter, r *http.Request) {

    // Get the query parameter for u
    uParam := r.FormValue("u")

    if uParam == "" {
        content, err := ioutil.ReadFile("./templates/splash.html")
        if err != nil {
            http.Error(rw, err.Error(), http.StatusNotFound)
            return
        }
        fmt.Fprintf(rw, string(content))
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

    // Create a template for the result
    fp := path.Join("templates", "index.html")
    tmpl, err := template.ParseFiles(fp)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
        return
    }

    // Serve the template
    if err := tmpl.ExecuteTemplate(rw, "filename", filename); err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}

// StaticHandler handles the static assets of the app.
func StaticHandler(rw http.ResponseWriter, r *http.Request) {

    // Only make .json and .css routes public
    if (strings.HasSuffix(r.URL.Path, ".json") || 
          strings.HasSuffix(r.URL.Path, ".css")  || 
          strings.HasSuffix(r.URL.Path, ".otf")  || 
          strings.HasSuffix(r.URL.Path, ".ico")) {
        http.ServeFile(rw, r, r.URL.Path[1:])
        return
    }
    http.NotFound(rw, r)
}