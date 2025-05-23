package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Page struct {
	Paragraph string
}

func Handler(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles("page.html"))
	p := Page{
		Paragraph: time.Now().String(),
	}
	err := templ.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
