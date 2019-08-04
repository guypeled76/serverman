package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

func main() {

	port := 8666

	fmt.Printf("Serverman running at port %v.\n", port)

	http.HandleFunc("/update/", handleUpdate)
	http.HandleFunc("/query/", handleQuery)
	http.HandleFunc("/", handleUnknown)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != nil {
		panic(err)
	}
}

func handleUnknown(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("Could not find a valid end point.\n"))
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	message := "===============================QUERY====================================\n"
	message += describeRequest(r)
	message += writeResponse(w, &struct {
		Status string  `json:"status"`
		Load   float32 `json:"load"`
	}{
		"ok",
		random.Float32(),
	})
	fmt.Print(message)
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	message := "===============================UPDATE===================================\n"
	message += describeRequest(r)
	message += writeResponse(w, &struct {
		Status string `json:"status"`
	}{
		"ok",
	})
	fmt.Print(message)
}

func writeResponse(w http.ResponseWriter, r interface{}) string {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	res, _ := json.Marshal(&r)
	w.Write(res)
	return fmt.Sprintf("response:%v\n", string(res))
}

func describeRequest(r *http.Request) string {

	var request []string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	if r.Method == "POST" {
		bodyBuffer, _ := ioutil.ReadAll(r.Body)
		request = append(request, fmt.Sprintf("body: \n%v\n", string(bodyBuffer)))
	}

	return strings.Join(request, "\n") + "\n"
}
