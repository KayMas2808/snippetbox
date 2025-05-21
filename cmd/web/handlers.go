package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

//
// Handler Functions
//

// home handles requests to "/"
// define a function as a func against the *application struct
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Makes it not act like a catch-all, and only work with URL "/"
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Initialize a slice containing the paths to the two files. It's important
	// to note that the file containing our base template must be the first
	// file in the slice.
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
	}

	// Use the template.ParseFiles() function to read the template file into a
	// template set. If there's an error, we log the detailed error message and use
	// the http.Error() function to send a generic 500 Internal Server Error
	// response to the user.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal server error", 500)
		return
	}

	// We then used the Execute() method on the template set to write the
	// template content as the response body. The last parameter to Execute()
	// represents any dynamic data that we want to pass in, which for now we'll
	// leave as nil.
	// now, we use the ExecuteTemplate() method to write the content of the "base"
	// template as the response body.
	// execute template tell go that we want to respond using the base template
	// which in turn invokes title and main templates
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal server error", 500)
	}
}

// snippetView handles requests to "/snippet/view"
func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id")) // get id parameter from url
	// and convert it to integer
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "display snippets for id = %d", id)
}

// snippetCreate handles POST requests to "/snippet/create"
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// .Set() function can set certain rules for the header, like the content type.
		// If we don't specify the content type, Go calls http.DetectContentType() automatically.
		// However, it can't differentiate JSON from plain text — so we explicitly define it.

		// Inform the client that only POST is allowed
		w.Header().Set("Allow", http.MethodPost) // uses constant instead of string
		// w.Header().Set("Allow", "POST")       // string alternative

		// If WriteHeader is not called, Go automatically sends a 200 OK.
		// It must be called before any call to w.Write.
		w.WriteHeader(http.StatusMethodNotAllowed) // using constant
		// w.WriteHeader(405)                      // integer alternative

		w.Write([]byte("Invalid http method. POST only allowed."))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=2000") // max seconds a response can be cached
	w.Write([]byte("{'name':'sam'}"))
}
