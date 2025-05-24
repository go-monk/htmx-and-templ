package main

import (
	"fmt"
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
	data := Page{
		FullReloadTime: time.Now().Format(time.DateTime),
	}
	PageLayout(data).Render(r.Context(), w)
}
