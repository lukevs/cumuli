// cumuli - A followings visualizer for SoundCloud

// To do:
//  - Add layout-based templating
//  - Add http-based error handling universally
//  - Add dynamic circle sizing
//  - Add logging
//  - Add the About page
//  - Link the JS

package main 

import (
    "log"
    "net/http"
)

func main() {
    // Routes
    http.HandleFunc("/", MainHandler)
    http.HandleFunc("/static/", StaticHandler)

    log.Println("Running on port :8080")
    http.ListenAndServe(":8080", nil)
}