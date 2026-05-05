package main

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func writeJSONError(w http.ResponseWriter, status int, err string, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   err,
		Message: msg,
	})
}

func hello(w http.ResponseWriter, req *http.Request) {
	// http.Error(w, "invalid request", http.StatusBadRequest) -> 400 Bad Request + Content-Type: text/plain
	writeJSONError(w, http.StatusNotFound, "error", "user not found")
}

func headers(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := make(map[string]any)
	for name, headers := range req.Header {
		for _, h := range headers {
			response[name] = h
		}
	}
	json.NewEncoder(w).Encode(response)

	// var apiResp map[string]any
	// if err := json.NewDecoder(res.Body).Decode(&apiResp); err != nil {
	// 	fmt.Println(err)
	// }
	// json.NewEncoder(w).Encode(apiResp)
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/hello/", hello)
	http.HandleFunc("/headers", headers)
	http.ListenAndServe(":8091", nil)
}
