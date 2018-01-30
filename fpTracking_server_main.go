//Website for the differents response types : http://www.alexedwards.net/blog/golang-response-snippets
//Website for the web sessions : http://www.gorillatoolkit.org/pkg/sessions
//Website for the ECMA6 : https://www.wanadev.fr/21-introduction-a-ecmascript-6-le-javascript-de-demain/

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
)

const TEMPLATES_FOLDER = "templates"

func handledFunctions() {
	http.HandleFunc("/",indexHandler)
	http.HandleFunc("/test-post/", testPostHandler)
	http.HandleFunc("/tracking-parallel/",trackingParallelHandler)
}

func main() {
	log.Println("Starting HTTP server on http://localhost:8080")

	// subscribe to SIGINT signals
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	srv := &http.Server{Addr: ":8080", Handler: http.DefaultServeMux}
	go func() {
		<-quit
		log.Println("Shutting down server...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("could not shutdown: %v", err)
		}
	}()

	//Handled functions
	handledFunctions()

	//Handled static files (like CSS and JS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	//Start the server and listen to requests
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}


	log.Println("Server stopped")
}