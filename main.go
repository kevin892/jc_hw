package main

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

//Stats data
type stats struct {
	Total   int     `json:"total"`
	Average float64 `json:"average"`
}

var finishedPasswords []string
var requestStats stats
var totalRequestTime float64

func hashPassword(password string) {
	hashedPassword := sha512.Sum512([]byte(password))
	encodePassword(hashedPassword)
}

func encodePassword(password [64]byte) {
	encodedPassword := base64.StdEncoding.EncodeToString(password[:])
	finishedPasswords = append(finishedPasswords, encodedPassword)
}

func addToCounter(duration int64) {
	totalRequestTime += float64(duration)
	requestStats.Total++
	requestStats.Average = totalRequestTime / float64(requestStats.Total)
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.String() {
	case "/hash":
		startTime := time.Now()
		r.ParseForm()
		if r.Form.Get("password") != "" {
			password := r.Form.Get("password")
			hashPassword(password)
			duration := time.Since(startTime).Microseconds()
			time.Sleep(5 * time.Second)
			fmt.Fprintf(w, finishedPasswords[requestStats.Total])
			addToCounter(duration)
		} else {
			fmt.Fprintf(w, "Error - 'password' not found!")
		}
	case "/shutdown":
		log.Fatalf("%s\n", "Shutting down")
	case "/stats":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(requestStats)
	default:
		return
	}
}

func handleRequests() {
	http.HandleFunc("/hash", requestHandler)
	http.HandleFunc("/shutdown", requestHandler)
	http.HandleFunc("/stats", requestHandler)
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleRequests()
}
