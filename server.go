package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

type User struct {
	Name string `json:"name"`
}

type NextActivityPage struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Assignments []Assignment
}

type Assignment struct {
	Title string
}

type Page struct {
	Title               string `json:"headline"`
	Description         string `json:"intro"`
	TitleLength         int    `json:"length_of_title"`
	HasAssignments      bool   `json:"has_assignments"`
	NumberOfAssignments int    `json:"number_of_assignments"`
}

func Adapt(nextPage NextActivityPage) Page {
	var adapted Page
	adapted.Title = nextPage.Title
	adapted.Description = nextPage.Description
	adapted.TitleLength = len(adapted.Title)
	adapted.HasAssignments = len(nextPage.Assignments) > 0
	adapted.NumberOfAssignments = len(nextPage.Assignments)

	return adapted
}

func main() {

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		peter := User{
			Name: "John",
		}

		json.NewEncoder(w).Encode(peter)
	})

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	})

	r.HandleFunc("/activityPage/{pageId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		requestURL := "https://web-next-api-dev.azurewebsites.net/api/ActivityPage/Get/" + vars["pageId"]
		res, err := http.Get(requestURL)
		if err != nil {
			fmt.Printf("error making http request: %s\n", err)
			os.Exit(1)
		}

		var nextPage NextActivityPage
		json.NewDecoder(res.Body).Decode(&nextPage)

		adapted := Adapt(nextPage)

		json.NewEncoder(w).Encode(adapted)
	})

	http.ListenAndServe(":8080", r)
}
