package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func main() {
	temp, err := template.ParseGlob("templates/*.html")
	if err != nil {
		fmt.Println(fmt.Sprint("ERREUR => %s", err.Error()))
		return
	}

	type Eleve struct {
		Name string
		Sexe int
		Age  int
	}

	type Promo struct {
		Title    string
		Filiere  string
		Level    int
		Students []Eleve
		Nb       int
	}

	http.HandleFunc("/promo", func(w http.ResponseWriter, r *http.Request) {
		data := Promo{"Mentors", "Informatique", 5, []Eleve{{"Cyrielle Rodrigues", 0, 22}, {"Kheir-eddine Mederreg", 1, 22}, {"Alan Philipiert", 1, 26}}, 3}

		temp.ExecuteTemplate(w, "promo", data)
	})

	RootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(RootDoc + "/asset/"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))

	http.ListenAndServe("localhost:8080", nil)
}
