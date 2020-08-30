package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"
)

var serverRequests int
var averageTime float64
var requestTimes float64
var finishedPasswords []string

func hashPassword(password string) {
	hashedPassword := sha512.Sum512([]byte(password))
	encodePassword(hashedPassword)
}

func encodePassword(password [64]byte) {
	encodedPassword := base64.StdEncoding.EncodeToString(password[:])
	finishedPasswords = append(finishedPasswords, encodedPassword)
}

func addToCounter(duration int64) {
	requestTimes = requestTimes + float64(duration)
	serverRequests++
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.String() {
	case "/hash":
		startTime := time.Now()
		r.ParseForm()
		if r.Form.Get("password") == "" {
			fmt.Fprintf(w, "Error - 'password' not found!")
		} else {
			password := r.Form.Get("password")
			hashPassword(password)
			duration := time.Since(startTime).Microseconds()
			time.Sleep(5 * time.Second)
			fmt.Fprintf(w, finishedPasswords[serverRequests])
			addToCounter(duration)
		}
	case "/shutdown":
		log.Fatalf("%s\n", "Shutting down")
	case "/stats":
		if serverRequests == 0 {
			fmt.Fprintf(w, "average: %v\n", 0)
		} else {
			fmt.Fprintf(w, "average: %v\n", requestTimes/float64(serverRequests))
		}
		fmt.Fprintf(w, "total: %v\n", serverRequests)
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
	serverRequests = 0
	averageTime = 0
	handleRequests()
}

// curl -w '\n' -d "password=angryMonkey" -X POST http://localhost:8080/hash
