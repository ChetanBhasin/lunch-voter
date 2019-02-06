package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/user"
	"strings"

	"github.com/manifoldco/promptui"
)

type lunchPlace struct {
	ID        string
	Name      string
	Distance  string
	PlaceType string
}

type votable struct {
	Places []lunchPlace
}

func selectPlaces(places []lunchPlace) lunchPlace {
	templates := promptui.SelectTemplates{
		Active:   `✔ {{ .Name | cyan | bold }}`,
		Inactive: `   {{ .Name | cyan }}`,
		Selected: `{{ "✔" | green | bold }} {{ "Recipe" | bold }}: {{ .Name | cyan }}`,
		Details: `
		--------- More Information ----------
		{{ "Name:" | faint }}	{{ .Name }}
		{{ "Distance:" | faint }}	{{ .Distance }}
		{{ "Type:" | faint }}	{{ .PlaceType }}`,
	}

	list := promptui.Select{
		Label:     "Select a place",
		Items:     places,
		Templates: &templates,
		Searcher: func(input string, index int) bool {
			return strings.Contains(places[index].Name, input)
		},
		StartInSearchMode: true,
	}

	index, _, err := list.Run()

	if err != nil {
		panic(err)
	}

	return places[index]
}

func ask(label string, defaultCase bool) bool {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	result, err := prompt.Run()

	if err != nil {
		return defaultCase
	}

	return result == "y"
}

func getUser() string {
	validate := func(input string) error {
		if len(input) < 3 {
			return errors.New("Username must have more than 3 characters")
		}
		return nil
	}

	var username string
	u, err := user.Current()
	if err == nil {
		username = u.Username
	}

	prompt := promptui.Prompt{
		Label:    "Username",
		Validate: validate,
		Default:  username,
	}

	result, err := prompt.Run()

	if err != nil {
		panic(err)
	}

	return result
}

func outputVotes() {
	files, err := ioutil.ReadDir("./.results")
	if err != nil {
		panic(err)
	}
	votes := make([]votable, 0)
	for _, f := range files {
		if !f.IsDir() {
			readable, err := ioutil.ReadFile("./.results/" + f.Name())
			if err != nil {
				panic(err)
			}
			var result votable
			err = json.Unmarshal(readable, &result)
			if err != nil {
				panic(err)
			}
			votes = append(votes, result)
		}
	}

	count := make(map[string]int)
	for _, vote := range votes {
		for _, place := range vote.Places {
			count[place.Name]++
		}
	}

	for place, number := range count {
		fmt.Println(number, " votes for ", place)
	}
}

func main() {

	if ask("Check resulst? Or ignore and proceed to voting?", false) {
		outputVotes()
		return
	}
	user := getUser()
	inFile, _ := ioutil.ReadFile("places.json")
	var data votable
	err := json.Unmarshal(inFile, &data)
	if err != nil {
		panic(err)
	}

	myPlaces := make([]lunchPlace, 0)

	voting := true

	for voting {
		myPlaces = append(myPlaces, selectPlaces(data.Places))
		voting = ask("Vote for another one?", false)
	}

	jsonOutput, _ := json.Marshal(votable{
		Places: myPlaces,
	})
	lineEnding := []byte("\n")

	unixOutput := bytes.Join([][]byte{jsonOutput, lineEnding}, []byte{})

	err = ioutil.WriteFile("./.results/"+user+".json", unixOutput, 0644)
	if err != nil {
		panic(err)
	}

}
