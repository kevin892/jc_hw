package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"
)

var serverRequests int
var averageTime float32
var reqTimes []float32
var allTimes float64

// func getTime(time) {

// }

func requestHandler(w http.ResponseWriter, r *http.Request) {
	srv := &http.Server{Addr: ":8080"}
	switch r.URL.String() {
	case "/hash":
		startTime := time.Now()
		r.ParseForm()
		fmt.Println(r.URL)
		password := r.Form.Get("password")
		hashedPassword := sha512.Sum512([]byte(password))
		encodedPassword := base64.StdEncoding.EncodeToString(hashedPassword[:])
		time.Sleep(5 * time.Second)
		fmt.Fprintf(w, encodedPassword)
		duration := time.Since(startTime).Milliseconds()
		allTimes = allTimes + float64(duration)
		serverRequests++
	case "/shutdown":
		srv.Close()
		fmt.Fprintf(w, "BYE FOR NOW")
	case "/stats":
		// allTimes = allTimes/float64(serverRequests)
		fmt.Fprintf(os.Stderr, "Need some stats?")
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
