package handlers

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles("./handlers/index.html"))

func MainPageHandler(w http.ResponseWriter, req *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
