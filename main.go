package main

import (
	"fmt"
	"net/http"
)

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	result := 10 + 5
	fmt.Fprintf(w, "Result: %d", result)
}

func main() {
	http.HandleFunc("/calculate", calculateHandler)

	port := 8989
	fmt.Printf("Server working on %d port...\n", port)

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
