package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
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

	type Visit struct {
		Number int
		Even   bool
	}

	type DataForm struct {
		FirstName string
		LastName  string
		BirthDate string
		Sexe      string
		Promo     string
		Error     bool
	}

	FormData := DataForm{Error: true}
	DisplayData := &FormData

	ChangeVisits := Visit{0, true}

	ListEleves := []Eleve{{"Cyrielle Rodrigues", 0, 22}, {"Kheir-eddine Mederreg", 3, 22}, {"Alan Philipiert", 1, 26}}

	http.HandleFunc("/promo", func(w http.ResponseWriter, r *http.Request) {
		data := Promo{"Mentors", "Informatique", 5, ListEleves, 3}

		temp.ExecuteTemplate(w, "promo", data)
	})

	http.HandleFunc("/change", func(w http.ResponseWriter, r *http.Request) {
		ChangeVisits.Number += 1
		ChangeVisits.Even = ChangeVisits.Number%2 == 0

		temp.ExecuteTemplate(w, "change", ChangeVisits)
	})

	http.HandleFunc("/user/init", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "form", 0)
	})

	http.HandleFunc("/user/treatment", func(w http.ResponseWriter, r *http.Request) {
		*DisplayData = DataForm{r.FormValue("firstname"), r.FormValue("lastname"), r.FormValue("birth"), r.FormValue("sexe"), r.FormValue("promo"), false}
		checkLastNameValue, _ := regexp.MatchString("^[a-zA-Z]{1,64}$", FormData.LastName)
		checkFirstNameValue, _ := regexp.MatchString("^[a-zA-Z]{1,64}$", FormData.FirstName)
		if !checkLastNameValue || !checkFirstNameValue {
			FormData.Error = true
		}

		http.Redirect(w, r, "/user/display", http.StatusSeeOther)
	})

	http.HandleFunc("/user/display", func(w http.ResponseWriter, r *http.Request) {
		if FormData.Error {
			http.Redirect(w, r, "/user/init", http.StatusForbidden)
		} else {
			temp.ExecuteTemplate(w, "display", FormData)
		}
	})

	RootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(RootDoc + "/asset/"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))

	http.ListenAndServe("localhost:8080", nil)
}
