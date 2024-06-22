package product

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (h *Handler) CountViewOnProductListingView(w http.ResponseWriter, r *http.Request) {
	// Log request method and path
	log.Printf("Received %s request for path: %s\n", r.Method, r.URL.Path)

	// Log request headers
	log.Println("Request headers:")
	for key, values := range r.Header {
		for _, value := range values {
			log.Printf("%s: %s\n", key, value)
		}
	}

	// Log request body (if present)
	if r.Body != nil {
		defer r.Body.Close()
		var requestData map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			log.Println("Error decoding request body:", err)
		} else {
			log.Println("Request body:")
			for key, value := range requestData {
				log.Printf("%s: %v\n", key, value)
			}
		}
	}

	// Your logic to handle the batch update goes here

	// Send response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Batch update completed successfully")
}
