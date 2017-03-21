package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/palicao/go-device-detection/lib"
)

func main() {
	address := ":8080"
	fmt.Printf("Listening to address %s\n", address)
	http.HandleFunc("/", handler)
	http.ListenAndServe(address, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	ua := r.UserAgent()
	parsed, _ := lib.Detect(ua)
	fmt.Fprintf(w, "%#v\n", parsed)

	elapsed := time.Since(start)
	fmt.Fprintf(w, "Detection took %s\n", elapsed)
}
