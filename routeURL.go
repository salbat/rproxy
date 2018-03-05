package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func NewReverseProxy(target string) *httputil.ReverseProxy {
	return httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   target,
	})
}

func Handle(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("request:", r.RemoteAddr, "want", r.RequestURI)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
		p.ServeHTTP(w, r)
	}
}

var Config map[string]interface{}

func main() {
	log.Print("init logger")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Print("load configuration")
	f, err := os.Open("./conf.json")
	if err != nil {
		log.Println(err.Error())
		return
	}
	if err := json.NewDecoder(f).Decode(&Config); err != nil {
		log.Println(err.Error())
		return
	}

	log.Print("init routes")
	for Path, Target := range Config["routes"].(map[string]interface{}) {
		if Path != "#" {
			http.HandleFunc(Path, Handle(NewReverseProxy(Target.(string))))
			log.Printf("%s > %s", Path, Target)
		}
	}

	Address := fmt.Sprintf("%s:%s", Config["ip"], Config["port"])
	log.Print("start listening on " + Address)
	http.ListenAndServe(Address, nil)
}
