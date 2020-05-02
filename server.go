package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

// StarWarsBaseURL is the base endpoint for swapi
const StarWarsBaseURL = "https://swapi.dev/api/"

// Planet describes a planet in Star Wars
type Planet struct {
	Name       string `json:"name"`
	Population string `json:"population"`
	Terrain    string `json:"terrain"`
}

// Character describes a character in Star Wars
type Character struct {
	Name         string `json:"name"`
	HomeWorldURL string `json:"homeworld"`
	HomeWorld    Planet
}

// AllCharacters describes a list of characters in Star Wars
type AllCharacters struct {
	Characters []Character `json:"results"`
}

// PageVariables describes the page-level variables for template
type PageVariables struct {
	PageTitle  string
	Characters []Character
}

// getHomeWorld retreives Planet information for a Character and updates the Character
func (character *Character) getHomeWorld(res http.ResponseWriter) *Character {
	var response *http.Response
	var err error
	if response, err = http.Get(character.HomeWorldURL); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Print("Error when requesting /api/planets/:id")
	}

	var bytes []byte
	if bytes, err = ioutil.ReadAll(response.Body); err != nil {
		log.Print("Error reading response body ", err)
	}

	if err := json.Unmarshal(bytes, &character.HomeWorld); err != nil {
		log.Print("Error parsing HomeWorld JSON ", err)
	}

	return character
}

// HomeHandler handles the "/" route and displays character information from Star Wars
func homeHandler(res http.ResponseWriter, req *http.Request) {
	var response *http.Response
	var err error
	if response, err = http.Get(StarWarsBaseURL + "people"); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Print("Error when requesting /api/people")
	}

	var bytes []byte
	if bytes, err = ioutil.ReadAll(response.Body); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Print("Error occurred while reading res body: ", err)
	}

	var allCharacters AllCharacters
	if err := json.Unmarshal(bytes, &allCharacters); err != nil {
		log.Print("Error parsing JSON ", err)
	}

	var updatedCharacters []Character
	for _, character := range allCharacters.Characters {
		updatedCharacters = append(updatedCharacters, *character.getHomeWorld(res))
	}

	pageVariables := PageVariables{
		PageTitle:  "Star Wars Characters",
		Characters: updatedCharacters,
	}

	t, err := template.ParseFiles("index.html")

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		log.Print("Template parsing error:", err)
	}

	err = t.Execute(res, pageVariables)
}

func main() {
	http.HandleFunc("/", homeHandler)
	fmt.Println("Server listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
