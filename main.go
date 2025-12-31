package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

var tmpl = template.Must(template.ParseFiles("template/index.html"))
var notesTmp = template.Must(template.ParseFiles("template/notes.html"))

type Notes struct {
	Id int
	Title string
	Note string
}
var notes []Notes

func main() {

	loadJson()

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
			Id: len(notes) + 1, 
			Title: title,
			Note: note,
		})

		saveJson()

		tmpl.Execute(w, notes)

		



	})

	http.HandleFunc("/show-note", func(w http.ResponseWriter, r *http.Request){
		
		if r.Method != http.MethodPost{
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		r.ParseForm()

		notesTmp.Execute(w, notes)
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request){
		if r.Method != http.MethodPost{
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		r.ParseForm()
		idstr := r.FormValue("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			fmt.Println("Invalid ID:", err)
        	http.Redirect(w, r, "/", http.StatusSeeOther)
        	return
		}

		for i, note := range notes {
			if note.Id == id {
				notes = append(notes[:1], notes[i+1:]... )
				break
			}
		}

		saveJson()

		http.Redirect(w, r, "/", http.StatusSeeOther)

		
	})

	fmt.Println("Server is running in http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}


func saveJson() {
	file, err := os.Create("data/notes.json")
	if err != nil {
		fmt.Println("Error creating file", err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")

	err = encoder.Encode(notes)
	if err != nil {
		fmt.Println("error encoding json", err)
	}
}

func loadJson() {
	file, err := os.Open("data/notes.json")
	if err != nil {
		fmt.Println("No existing note file", err)
		return
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&notes)
	if err != nil {
		fmt.Println("Error decoding json file:", err)
	}

}

