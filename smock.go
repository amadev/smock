package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Response struct {
	Code uint
	Data string
}

type ConfRequest struct {
	Path string
	Code uint
	Data string
}

type Conf map[string]Response

var CNF Conf

const CONFIG_PATH = "/_cnf"

func cnf(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		log.Println("Listing configuration", req.URL)
		u, _ := json.Marshal(CNF)
		fmt.Fprint(w, string(u))
	} else if req.Method == http.MethodPost {
		c := &ConfRequest{}
		err := json.NewDecoder(req.Body).Decode(c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Setting configuration, data: %+v\n", c)
		CNF[c.Path] = Response{Code: c.Code, Data: c.Data}
		fmt.Fprint(w, "OK")
	} else if req.Method == http.MethodDelete {
		c := &ConfRequest{}
		err := json.NewDecoder(req.Body).Decode(c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Deleting configuration, data: %+v\n", c)
		delete(CNF, c.Path)
		fmt.Fprint(w, "OK")
	}
}

func PathHandler(handler http.Handler) http.Handler {

	log.Print("starting handler mux")

	CNF = make(map[string]Response)
	CNF["/success"] = Response{Code: 200, Data: "{\"result\": \"OK\"}"}
	CNF["/fail"] = Response{Code: 500, Data: "{\"result\": \"FAIL\"}"}
	CNF["/bad"] = Response{Code: 400, Data: "{\"result\": \"BAD\"}"}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		log.Printf("Got request, url: %s, method: %s\n", path, req.Method)

		if path == CONFIG_PATH {
			handler.ServeHTTP(w, req)
		} else if val, ok := CNF[path]; ok {
			log.Printf("Found response, path: %s, data: %+v\n", path, val)

			w.WriteHeader(int(val.Code))
			w.Write([]byte(val.Data))
		} else if val, ok := CNF["/default"]; ok {
			log.Printf("Found default response, path: %s, data: %+v\n", req.URL, val)

			w.WriteHeader(int(val.Code))
			w.Write([]byte(val.Data))
		} else {
			log.Printf("Standard response")
			handler.ServeHTTP(w, req)
		}

		log.Printf("Processed request, url: %s, method: %s\n", req.URL, req.Method)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(CONFIG_PATH, cnf)

	WrappedMux := PathHandler(mux)
	http.ListenAndServe(":8080", WrappedMux)
}
