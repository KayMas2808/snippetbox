package main

import (
	"log"
	"net/http"
)

//
// Entry Point
//

func main() {
	mux := http.NewServeMux()

	// Route registrations
	mux.HandleFunc("/", home)
	// "/" is treated like a catch-all -> all requests at any URL not routed elsewhere
	// will be handled by the home function.

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/snippet/create", snippetCreate) // fixed URL path
	mux.HandleFunc("/snippet/view", snippetView)     // fixed URL path

	// ServeMux notes:
	// - It always gives precedence to longer URL paths.
	// - Subtree paths like "/snippet/" are handled by prefix matching.
	// - Paths like "/snippet" are redirected to "/snippet/".
	// - Auto path cleaning: e.g., "/snippet/.//." becomes "/snippet/".

	log.Println("Starting servemux on :8080")

	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
