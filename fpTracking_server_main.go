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
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const TEMPLATES_FOLDER = "templates"

//Global variables
var progressInformationSession map[string]*progressInformationStruct
var (
	key = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)


func main() {
	log.Println("Starting HTTP server on http://localhost:8080")

	//Create the router and handle the URLs
	r := mux.NewRouter()
	r.HandleFunc("/",indexHandler)
	r.HandleFunc("/test-post", testPostHandler)
	r.HandleFunc("/tracking-parallel",trackingParallelHandler)
	r.HandleFunc("/check-progression",checkProgressionHandler)

	//Handled static files (like CSS and JS)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))


	// subscribe to SIGINT signals
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	srv := &http.Server{
		Handler: r,
		Addr: ":8080",
	}
	go func() {
		<-quit
		log.Println("Shutting down server...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("could not shutdown: %v", err)
		}
	}()


	//Initialization of the global variables
	progressInformationSession = make(map[string]*progressInformationStruct)
	log.Println("progressInformationSession :",progressInformationSession)

	//Start the server and listen to requests
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}


	log.Println("Server stopped")
}