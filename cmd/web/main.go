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
	// new http.Server struct to use our custom error logger for Go's http server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting servemux on %s", *addr)

	err := srv.ListenAndServe()
	// we could've used os.Getenv("var name") to get from environment variables
	// but that dont have a default value, and return type is always string
	errorLog.Fatal(err)
}
