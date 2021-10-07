package main

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // retrieve key from the request
	key := vars["key"]

	value, err := io.ReadAll(r.Body) // the request body has our value
	defer r.Body.Close()

	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	err = Put(key, string(value)) // Store the value as a string
	if err != nil {               // If we have an error, report it
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated) // All good! Return StatusCreated
}

func main() {
	r := mux.NewRouter()

	// Register keyValuePutHandler as the handler function for PUT
	// requests matching "/v1/{key}"
	r.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", r))
}
