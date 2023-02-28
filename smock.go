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

func cnf(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		u, _ := json.Marshal(CNF)
		fmt.Fprint(w, string(u))
	} else if req.Method == http.MethodPost {
		c := &ConfRequest{}
		err := json.NewDecoder(req.Body).Decode(c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		CNF[c.Path] = Response{Code: c.Code, Data: c.Data}
		fmt.Fprint(w, "OK")
	} else if req.Method == http.MethodDelete {
		c := &ConfRequest{}
		err := json.NewDecoder(req.Body).Decode(c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
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
		log.Println("Got request", req.URL)
		if val, ok := CNF[req.URL.Path]; ok {
			w.WriteHeader(int(val.Code))
			w.Write([]byte(val.Data))
		} else {
			handler.ServeHTTP(w, req)

		}
		log.Println("Processed request", req.URL)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/_cnf", cnf)

	WrappedMux := PathHandler(mux)
	http.ListenAndServe(":8080", WrappedMux)
}
