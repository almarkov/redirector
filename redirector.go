package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

// Route - redirection route
type Route struct {
	name string
	maps map[string]string
}

// Config - application configuration
type Config struct {
	address string
	routes  []Route
}

var config Config

func main() {
	config = ReadConfig()
	for _, element := range config.routes {
		if element.name != "" {
			http.HandleFunc(element.name, handler)
		}
	}
	log.Fatal(http.ListenAndServe(config.address, nil))
}

// ReadConfig - read aplication configuration
func ReadConfig() Config {

	file, err := os.Open("config")
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	var c Config
	var cr int
	c.routes = make([]Route, 10)
	cr = -1
	for _, eachline := range txtlines {
		if len(eachline) > 0 {
			line := strings.Split(eachline, " ")
			if line[0] == "address" {
				c.address = line[1]
			} else if line[0] == "route" {
				cr++
				c.routes[cr].maps = make(map[string]string)
				c.routes[cr].name = line[1]
			} else {
				c.routes[cr].maps[line[0]] = line[1]
			}
		}
	}

	return c
}

func getip(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded == "" {
		forwarded = r.RemoteAddr
	}
	host, _, err := net.SplitHostPort(forwarded)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ip fetch error: %v\n", err)
	}
	return host
}

func handler(w http.ResponseWriter, r *http.Request) {

	clientIP := getip(r)

	// russian ip for test
	if clientIP == "127.0.0.1" {
		clientIP = "37.232.169.197"
	}

	resp, err := http.Get("http://ip-api.com/json/" + clientIP)

	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
	}

	b, err := ioutil.ReadAll(resp.Body)

	var f interface{}

	err = json.Unmarshal(b, &f)
	m := f.(map[string]interface{})
	var s string
	for k, v := range m {
		if k == "countryCode" {
			s = v.(string)
		}
	}

	var redirectstr string
	for _, route := range config.routes {
		if route.name == r.URL.Path {
			redirectstr = route.maps[s]
			if redirectstr == "" {
				redirectstr = route.maps["DEFAULT"]
			}
		}
	}

	http.Redirect(w, r, redirectstr, 302)

}
