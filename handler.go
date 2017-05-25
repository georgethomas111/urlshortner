package urlshortner

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/georgethomas111/urlshortner/random"
	"github.com/gorilla/mux"
)

var hFactory *Factory
var rand *random.Random

// FIXME urlLength should be taken from the config
const urlLength = 9

// Redirects if the url is not returned from the store.
// TODO Logs can be used to save the access time.
// TODO The information can also be stored in the database.
// TODO Database implementation is preferred over implementation using logs
func redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tinyURL := vars["tinyurl"]
	longURL, err := hFactory.GetURL(tinyURL)
	if err != nil {
		//FIXME Establish the right method to log errors.
		http.NotFound(w, r)
	}
	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}

func add(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	longURL := vars["longurl"]
	// FIXME Adding default protocol.
	// TODO Accept the protocol from user instead of assuming.
	protocol := "https://"
	longURL = protocol + longURL

	var tinyURL string
	var err error
	var uniqueFound bool

	for uniqueFound == false {
		tinyURL, err = rand.GetRandomUrl(urlLength)
		if err != nil {
			// Log the error to know about the details.
			http.Error(w, "Could not create tinyURL", http.StatusInternalServerError)
		}

		err = hFactory.AddURL(tinyURL, longURL)
		if err != nil {
			// Log the error to know about the details.
			// Error while adding url try another url.
			continue
		}
		uniqueFound = true
	}

	resp := &struct {
		TinyURL string `json:"tinyurl"`
		LongURL string `json:"longurl"`
	}{
		tinyURL,
		longURL,
	}
	w.Header().Set("Content-Type", "json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(resp)
	// This might never happen.
	if err != nil {
		http.Error(w, "Encode error", http.StatusInternalServerError)
	}
}

// FIXME The arguments to the function Start should be taken in the config.
func Start(lAddr string, dbHost string, port int, db string, desDoc string) error {
	hFactory = NewFactory(dbHost, port, db, desDoc)
	var err error
	rand, err = random.NewRandom()
	if err != nil {
		return err
	}
	router := mux.NewRouter()
	router.HandleFunc("/{tinyurl}", redirect).Methods("GET")
	router.HandleFunc("/add/{longurl}", add).Methods("POST")
	fmt.Println("Welcome to url redirect. Listening on ", lAddr)
	http.ListenAndServe(lAddr, router)
	return nil
}
