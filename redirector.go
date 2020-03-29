package main

import (
	"encoding/json"
	"fmt"
	// "html/template"
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"os"
)

type IPInfo2 struct {
	country     string
	countryCode string
	query       string
}

type Route struct {
	name    string
	maps    map[string]string
}

type Config struct {
	address string
	routes  []Route
}

var config Config
func main() {
	config = ReadConfig()
	for _, element := range config.routes {
		if element.name != ""  {
			http.HandleFunc(element.name, handler)
		}
	}
	log.Fatal(http.ListenAndServe(config.address, nil))
}

func ReadConfig() Config{
	
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
 	var currrentroute int
 	c.routes = make([]Route, 10)
 	currrentroute = -1
	for _, eachline := range txtlines {
		if len(eachline) > 0 {
			line := strings.Split(eachline, " ")
			if line[0] == "address" {
				c.address = line[1]
			} else if line[0] == "route" {
				currrentroute += 1
				c.routes[currrentroute].maps = make(map[string]string)
				c.routes[currrentroute].name = line[1]
			} else {
				c.routes[currrentroute].maps[line[0]] = line[1]
			}
		}
	}

	return c
}

func getip(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func handler(w http.ResponseWriter, r *http.Request) {

	clientIP := getip(r)
	res := strings.Split(clientIP, ":")
	clientIP = res[0]

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
		if (route.name == r.URL.Path) {
			redirectstr = route.maps[s]
			if redirectstr == "" {
				redirectstr = route.maps["DEFAULT"]
			}
		}
	}

	http.Redirect(w, r, redirectstr, 302)

}
