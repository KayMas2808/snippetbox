package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// application struct to make loggers available in all files in this package
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

//
// Entry Point
//

func main() {

	// new command line flag, default value of 8080, usage text telling what it does.
	// couldve been flat.Int, flag.Bool etc
	addr := flag.String("addr", ":8080", "HTTP network address (defualt \":8080\")")

	// parse the flag - reads it in command line and assigns value to addr.
	// call this before using addr
	flag.Parse()

	// new logger to write information messages
	// 3 params: destination of logs(os.Stdout),
	//	 		 string prefix (INFO then tab),
	//	 		 flags for additional info (date and time - joined using bitwise OR)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	// Lshortfile adds relevant file name and number.

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := http.NewServeMux()

	// Route registrations
	mux.HandleFunc("/", app.home)
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

	mux.HandleFunc("/snippet/create", app.snippetCreate) // fixed URL path
	mux.HandleFunc("/snippet/view", app.snippetView)     // fixed URL path

	// ServeMux notes:
	// - It always gives precedence to longer URL paths.
	// - Subtree paths like "/snippet/" are handled by prefix matching.
	// - Paths like "/snippet" are redirected to "/snippet/".
	// - Auto path cleaning: e.g., "/snippet/.//." becomes "/snippet/".

	// new http.Server struct to use our custom error logger for Go's http server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting servemux on %s", *addr)

	err := srv.ListenAndServe()
	// we could've used os.Getenv("var name") to get from environment variables
	// but that dont have a default value, and return type is always string
	errorLog.Fatal(err)
}
