package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles("template/index.html"))

type Notes struct {
	Title string
	Note string
}
var notes []Notes

func main() {

	http.Handle("/static/", 
	http.StripPrefix("/static/", 
	http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		tmpl.Execute(w, notes)
	})

	http.HandleFunc("/add-note", func(w http.ResponseWriter, r *http.Request){

		if r.Method != http.MethodPost{
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		r.ParseForm()

		title := r.FormValue("title")
		note := r.FormValue("note")

		notes = append(notes, Notes{
			Title: title,
			Note: note,
		})

		tmpl.Execute(w, notes)

		



	})

	fmt.Println("Server is running in http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}