package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Record struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
	ID   string `json:"id"`
}

type recordHandlers struct {
	sync.Mutex
	store map[string]Record
}

func (h *recordHandlers) records(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func (h *recordHandlers) get(w http.ResponseWriter, r *http.Request) {
	records := make([]Record, len(h.store))

	// read the records from the memory store
	i := 0
	h.Lock()
	for _, record := range h.store {
		records[i] = record
		i++
	}
	h.Unlock()

	// convert to json
	jsonBytes, err := json.Marshal(records)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// send back the response
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (h *recordHandlers) getRandomRecord(w http.ResponseWriter, r *http.Request) {
	ids := make([]string, len(h.store))

	i := 0
	h.Lock()
	for id := range h.store {
		ids[i] = id
		i++
	}
	defer h.Unlock()

	var target string
	if len(ids) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if len(ids) == 1 {
		target = ids[0]
	} else {
		rand.Seed(time.Now().UnixNano())
		target = ids[rand.Intn(len(ids))]
	}

	w.Header().Add("location", fmt.Sprintf("/records/%s", target))
	w.WriteHeader(http.StatusFound)
}

func (h *recordHandlers) getRecord(w http.ResponseWriter, r *http.Request) {
	// parse the path
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if parts[2] == "random" {
		h.getRandomRecord(w, r)
		return
	}

	// read the records from the memory store
	h.Lock()
	record, ok := h.store[parts[2]]
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// convert to json
	jsonBytes, err := json.Marshal(record)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// send back the response
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *recordHandlers) post(w http.ResponseWriter, r *http.Request) {

	// read the body
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// check the content-type
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json' but got '%s'", ct)))
		return
	}

	// parse the body
	var record Record
	err = json.Unmarshal(bodyBytes, &record)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// generate a timestamp as ID
	record.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	h.Lock()
	h.store[record.ID] = record
	defer h.Unlock()
}

func newRecordHandlers() *recordHandlers {
	return &recordHandlers{
		// initialize an empty memory store
		store: map[string]Record{},
	}
}

type adminPortal struct {
	password string
}

func newAdminPortal() *adminPortal {
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		panic("required env var ADMIN_PASSWORD not set")
	}

	return &adminPortal{password: password}
}

func (a adminPortal) handler(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok || user != "admin" || pass != a.password {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - unauthorized"))
		return
	}

	w.Write([]byte("<html><h1>Super secret admin portal</h1></html>"))
}

func main() {
	admin := newAdminPortal()
	recordHandlers := newRecordHandlers()
	http.HandleFunc("/records", recordHandlers.records)
	http.HandleFunc("/records/", recordHandlers.getRecord)
	http.HandleFunc("/admin", admin.handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
