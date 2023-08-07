package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type InputData struct {
	A int `json:"a"`
	B int `json:"b"`
}

type OutputData struct {
	ResultA int `json:"a"`
	ResultB int `json:"b"`
}

func calculateFactorial(n int, ch chan<- int) {
	result := 1
	for i := 1; i <= n; i++ {
		result *= i
	}
	ch <- result
}

func calculateHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var inputData InputData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&inputData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if inputData.A < 0 || inputData.B < 0 {
		http.Error(w, `{"error":"Incorrect input"}`, http.StatusBadRequest)
		return
	}

	resultChA := make(chan int)
	resultChB := make(chan int)
	go calculateFactorial(inputData.A, resultChA)
	go calculateFactorial(inputData.B, resultChB)

	resultA, resultB := <-resultChA, <-resultChB

	outData := OutputData{
		ResultA: resultA,
		ResultB: resultB,
	}

	resultJSON, err := json.Marshal(outData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resultJSON)
}

func main() {
	router := httprouter.New()
	router.POST("/calculate", calculateHandler)

	port := 8989

	fmt.Printf("Server working on %d port ...\n", port)

	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
