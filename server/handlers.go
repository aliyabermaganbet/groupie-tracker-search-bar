package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

type Artist struct {
	Loc          []string
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	Creationdate int      `json:"creationDate"`
	Firstalbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	Concertdates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}
type Display struct {
	Name         string
	Image        string
	Member       []string
	Creationdate int
	Firstalbum   string
	Relate       map[string][]string
	Results      string
}
type Second struct {
	Index []struct {
		Id                int                 `json:"id"`
		Datesandlocations map[string][]string `json:"datesLocations"`
	}
}
type Third struct {
	Locations []string `json:"locations"`
}

func Homepage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	tmp, err := template.ParseFiles("template/html.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rspData, err1 := ChangeAPI("https://groupietrackers.herokuapp.com/api/artists")
	if err1 == 500 {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var artist []Artist
	jsonerr := json.Unmarshal(rspData, &artist)
	if jsonerr != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for i := 0; i < len(artist); i++ {
		rs, _ := ChangeAPI(artist[i].Locations)
		temp := Third{}
		jsonerr := json.Unmarshal(rs, &temp)
		if jsonerr != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		artist[i].Loc = temp.Locations
	}
	tmp.Execute(w, artist)
}

func Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	notfound, _ := template.ParseFiles("template/notfound.html")
	tmp, err := template.ParseFiles("template/test.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rspData, err1 := ChangeAPI("https://groupietrackers.herokuapp.com/api/artists")
	if err1 == 500 {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var artist []Artist
	jsonerr := json.Unmarshal(rspData, &artist)
	if jsonerr != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	flag := false
	var searchArtists []Artist
	if err := r.ParseForm(); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	opt, ok := r.Form["option"]
	if !ok {
		http.Error(w, "option field not found", http.StatusInternalServerError)
		return
	}
	search, ok := r.Form["search"]
	if !ok {
		http.Error(w, "search field not found", http.StatusInternalServerError)
		return
	}

	option := opt[0]
	inp := search[0]

	if option == "All" {
		for k := 0; k < len(artist); k++ { // here starts
			rs, _ := ChangeAPI(artist[k].Locations)
			temp := Third{}
			jsonerr := json.Unmarshal(rs, &temp)
			if jsonerr != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			stringOfLocations := strings.Join(temp.Locations, " ")
			if strings.Contains(stringOfLocations, inp) {
				searchArtists = append(searchArtists, artist[k])
				flag = true
				continue
			}
			for i := 0; i < len(artist[k].Members); i++ { // next step
				if strings.Contains(artist[k].Members[i], inp) || strings.Contains(artist[k].Name, inp) || strings.Contains(artist[k].Firstalbum, inp) || strings.Contains(strconv.Itoa(artist[k].Creationdate), inp) {
					searchArtists = append(searchArtists, artist[k])
					flag = true
					break
				}
			}
		}
		if flag {
			tmp.Execute(w, searchArtists)
		} else {
			notfound.Execute(w, nil)
		}
	} else if option == "Artist/band" {
		for _, el := range artist {
			if strings.Contains(el.Name, inp) {
				searchArtists = append(searchArtists, el)
				flag = true
				break
			}
		}
		if flag {
			tmp.Execute(w, searchArtists)
		} else {
			notfound.Execute(w, nil)
		}
	} else if option == "First album date" {
		for _, el := range artist {
			if strings.Contains(el.Firstalbum, inp) {
				searchArtists = append(searchArtists, el)
				flag = true
			}
		}
		if flag {
			tmp.Execute(w, searchArtists)
		} else {
			notfound.Execute(w, nil)
		}
	} else if option == "Creation date" {
		for _, el := range artist {
			if strings.Contains(strconv.Itoa(el.Creationdate), inp) {
				searchArtists = append(searchArtists, el)
				flag = true
			}
		}
		if flag {
			tmp.Execute(w, searchArtists)
		} else {
			notfound.Execute(w, nil)
		}
	} else if option == "Members" {
		for _, el := range artist {
			for _, k := range el.Members {
				if strings.Contains(k, inp) {
					searchArtists = append(searchArtists, el)
					flag = true
					break
				}
			}
		}
		if flag {
			tmp.Execute(w, searchArtists)
		} else {
			notfound.Execute(w, nil)
		}
	} else if option == "Locations" {
		for _, el := range artist { // here starts
			rs, _ := ChangeAPI(el.Locations)
			temp := Third{}
			jsonerr := json.Unmarshal(rs, &temp)
			if jsonerr != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			stringOfLocations := strings.Join(temp.Locations, " ")
			if strings.Contains(stringOfLocations, inp) {
				searchArtists = append(searchArtists, el)
				flag = true
				continue
			}
		}
		if flag {
			tmp.Execute(w, searchArtists)
		} else {
			notfound.Execute(w, nil)
		}
	}
}

func ChangeAPI(s string) ([]byte, int) {
	response, err := http.Get(s)
	if err != nil {
		return nil, 500
	}
	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, 500
	}
	return responseData, 0
}

func DisplayOutput(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	myid, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if myid > 52 || myid == 0 {
		http.NotFound(w, r)
		return
	}
	res, err := template.ParseFiles("template/result.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	responseData, err1 := ChangeAPI("https://groupietrackers.herokuapp.com/api/artists")
	if err1 == 500 {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var artist []Artist
	jsonerr := json.Unmarshal(responseData, &artist)
	if jsonerr != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rspData, err1 := ChangeAPI("https://groupietrackers.herokuapp.com/api/relation")
	if err1 == 500 {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var relate Second
	jsoner := json.Unmarshal(rspData, &relate)

	if jsoner != nil {
		fmt.Println(jsoner.Error())
	}
	myMap := make(map[string][]string)
	var members []string
	image := ""
	name := ""
	creationdate := 0
	firstalbum := ""
	for _, zn := range artist {
		for _, k := range relate.Index {
			if myid == zn.Id && myid == k.Id {
				name = zn.Name
				members = zn.Members
				image = zn.Image
				creationdate = zn.Creationdate
				firstalbum = zn.Firstalbum
				myMap = k.Datesandlocations
			}
		}
	}
	result := Display{
		Name:         name,
		Image:        image,
		Member:       members,
		Creationdate: creationdate,
		Firstalbum:   firstalbum,
		Relate:       myMap,
	}
	res.Execute(w, result)
}
