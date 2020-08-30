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
var allTimes float64

func requestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.String() {
	case "/hash":
		startTime := time.Now()
		r.ParseForm()
		password := r.Form.Get("password")
		hashedPassword := sha512.Sum512([]byte(password))
		encodedPassword := base64.StdEncoding.EncodeToString(hashedPassword[:])
		time.Sleep(5 * time.Second)
		fmt.Fprintf(w, encodedPassword)
		duration := time.Since(startTime).Milliseconds()
		allTimes = allTimes + float64(duration)
		serverRequests++
	case "/shutdown":
		// if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("%s\n", "Shutting down")
		// }
		return
	case "/stats":
		if serverRequests == 0 {
			fmt.Fprintf(w, "Duraction: %v\n", 0)
		} else {
			fmt.Fprintf(w, "Average Request Speed : %v\n", allTimes/float64(serverRequests))
		}
		fmt.Fprintf(w, "Server Requests: %v\n", serverRequests)
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
