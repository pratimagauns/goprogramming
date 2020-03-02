package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type customHandler struct{}

// ServeHTTP implements the http.Handler interface in the net/http package
func (h customHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// ParseForm will parse query string values and make r.Form available
	r.ParseForm()

	// r.Form is map of query string parameters
	// its' type is url.Values, which in turn is a map[string][]string
	queryMap := r.Form

	switch r.Method {
	case http.MethodGet:
		// Handle GET requests
		w.WriteHeader(http.StatusOK)
		fmt.Println("Processing... ** ", queryMap)

		records, err := ProcessDateRequest(queryMap)
		if err != "" {
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		w.Write([]byte(fmt.Sprintf("Records: %s", records)))
		return
	case http.MethodPost:
		// Handle POST requests
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			// Error occurred while parsing request body
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Query string values: %s\nBody posted: %s", queryMap, body)))
		return
	}

	// Other HTTP methods (eg PUT, PATCH, etc) are not handled by the above
	// so inform the client with appropriate status code
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func main() {

	// Create a mux for routing incoming requests
	m := http.NewServeMux()

	// All URLs will be handled by this function
	m.Handle("/fetch", customHandler{})

	// Create a server listening on port 8000
	s := &http.Server{
		Addr:    ":8000",
		Handler: m,
	}

	// Continue to process new requests until an error occurs
	log.Fatal(s.ListenAndServe())
}
