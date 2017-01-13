package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"bitbucket.com/ocono/device-detector/lib"
)

var browserRegexps = lib.ParseYml("./yml/browsers.yml")

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	ua := r.UserAgent()
	parsed := parse(ua)
	fmt.Fprintf(w, "%v", parsed)
	fmt.Printf("%v", parsed)

	elapsed := time.Since(start)
	fmt.Printf("Detection took %s\n", elapsed)
}

func parse(ua string) string {
	for _, br := range browserRegexps {
		//if !strings.Contains(br.Regex, "?!") {
		if br.Compiled != nil {
			found := br.Compiled.FindStringSubmatch(ua)

			if len(found) > 0 {
				browser := br.Name
				version := ""
				if br.Version[0:1] == "$" {
					part, _ := strconv.Atoi(br.Version[1:])
					version = found[part]
				} else {
					version = br.Version
				}
				return fmt.Sprintf("Browser: %s, Version: %s\n", browser, version)
			}
		}
	}
	return ""
}
