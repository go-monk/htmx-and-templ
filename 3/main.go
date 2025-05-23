package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", Handler)
	http.HandleFunc("/time", TimeHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Page struct {
	FullReloadTime string
}

func TimeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, time.Now().Format(time.DateTime))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles("page.html"))
	p := Page{
		FullReloadTime: time.Now().Format(time.DateTime),
	}
	if err := templ.Execute(w, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
