package main

import (
	"encoding/json"
	"fmt"
	// "html/template"
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

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("95.213.252.40:3060", nil))
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
	// clientIP := "37.232.169.197"

	resp, err := http.Get("http://ip-api.com/json/" + clientIP)

	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s\n", b)
	resp.Body.Close()
	fmt.Printf("decode error")
	fmt.Printf("%s\n", "there")
	// if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
	// 	resp.Body.Close()
	// 	fmt.Printf("decode error")
	// }
	// if err := json.Unmarshal([]byte(`{"country":"a","countryCode":"B","query":"ada"`), &result); err != nil {
	// 	fmt.Printf("%s\n", "unm err")
	// }
	var f interface{}
	// b1 := []byte(`{"countryCode":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)
	err = json.Unmarshal(b, &f)
	m := f.(map[string]interface{})
	var s string
	for k, v := range m {
		if k == "countryCode" {
			s = v.(string)
		}
	}
	fmt.Fprintf(w, "country code:%s\n", s)

}
