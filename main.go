package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type AllArtists struct {
	Id             int                 `json:"id"`
	Image          string              `json:"image"`
	Name           string              `json:"name"`
	Members        []string            `json:"members"`
	CreationDate   int                 `json:"creationDate"`
	FirstAlbum     string              `json:"firstAlbum"`
	Locations      []string            `json:"locations"`
	ConcertDates   []string            `json:"concertDates"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type Artist struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Locations struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Relation struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type Dates struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type DatesIndex struct {
	Index []Dates `json:"index"`
}

type LocationsIndex struct {
	Index []Locations `json:"index"`
}

type RelationsIndex struct {
	Index []Relation `json:"index"`
}

var (
	allArtistsData []AllArtists
	artistsData    []Artist
	datesData      DatesIndex
	locationsData  LocationsIndex
	relationsData  RelationsIndex
)

func main() {
	GetRelationsData()

	fs := http.FileServer(http.Dir("templates"))
	http.Handle("/templates/", http.StripPrefix("/templates/", fs))

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/artist", artistHandler)

	println("Server running")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Listen and Server", err)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	err := GetData()
	if err != nil {
		errors.New("Error by get data")
	}

	if r.URL.Path != "/" {
		http.Redirect(w, r, "templates/404.html", http.StatusFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := tmpl.Execute(w, allArtistsData); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, _ := strconv.Atoi(idStr)

	if r.URL.Path != "/artist" {
		http.Redirect(w, r, "templates/404.html", http.StatusFound)
		returnc
	}

	for _, artist := range allArtistsData {
		if artist.Id == id {

			tmpl, err := template.ParseFiles("templates/artist.html")
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}

			if err := tmpl.Execute(w, artist); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}
	}
}

func GetArtistsData() error {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return errors.New("Error")
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error")
	}
	json.Unmarshal(bytes, &artistsData)
	return nil
}

func GetDatesData() error {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/dates")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)

	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error")
	}

	json.Unmarshal(bytes, &datesData)
	return nil
}

func GetLocationsData() error {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)

	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error")
	}
	json.Unmarshal(bytes, &locationsData)
	return nil
}

func GetRelationsData() {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		log.Println(err.Error())
		return
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}
	json.Unmarshal(bytes, &relationsData)
}

func GetData() error {
	if len(allArtistsData) != 0 {
		return nil
	}
	err1 := GetArtistsData()
	err2 := GetLocationsData()
	err3 := GetDatesData()

	if err1 != nil || err2 != nil || err3 != nil {
		return errors.New("Error by get data artists, locations, dates")
	}
	for i := range artistsData {
		var temp AllArtists
		temp.Id = i + 1
		temp.Image = artistsData[i].Image
		temp.Name = artistsData[i].Name
		temp.Members = artistsData[i].Members
		temp.CreationDate = artistsData[i].CreationDate
		temp.FirstAlbum = artistsData[i].FirstAlbum
		temp.Locations = locationsData.Index[i].Locations
		temp.ConcertDates = datesData.Index[i].Dates
		temp.DatesLocations = relationsData.Index[i].DatesLocations
		allArtistsData = append(allArtistsData, temp)
	}
	return nil
}
