package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type ViewData struct {
	Woo []struct {
		Id          int      `json:"id"`
		Slug        string   `json:"slug"`
		DisplayName string   `json:"displayName"`
		FullName    string   `json:"fullName"`
		Species     string   `json:"species"`
		Sex         string   `json:"sex"`
		Quotes      []string `json:"quotes"`
		Sprite      string   `json:"sprite"`
		Background  string   `json:"background"`
	} `json:"woo"`
}

func loadAPI() ViewData {
	vd := ViewData{}

	url := "https://raw.githubusercontent.com/Mkheir13/Adventure/main/data/data.json"

	httpClient := http.Client{
		Timeout: time.Second * 2, // define timeout
	}

	//create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "API AT test <3")

	//make api call
	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	//parse response
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonErr := json.Unmarshal(body, &vd)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return vd
}

func main() {

	viewData := loadAPI()

	tmpol := template.Must(template.ParseFiles("../html/index.html"))

	cssFolder := http.FileServer(http.Dir("../css"))
	http.Handle("/css/", http.StripPrefix("/css/", cssFolder))

	imgFolder := http.FileServer(http.Dir("../img"))
	http.Handle("/img/", http.StripPrefix("/img/", imgFolder))

	jsFolder := http.FileServer(http.Dir("../js"))
	http.Handle("/js/", http.StripPrefix("/js/", jsFolder))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		search := r.FormValue("searchBar")
		if search != "" {
			filteredViewData := ViewData{}
			for _, adventure := range viewData.Woo {
				if strings.Contains(strings.ToLower(adventure.DisplayName), strings.ToLower(search)) || strings.Contains(strings.ToLower(adventure.FullName), strings.ToLower(search)) {
					filteredViewData.Woo = append(filteredViewData.Woo, adventure)
				}
			}
			tmpol.Execute(w, filteredViewData)
		} else {
			tmpol.Execute(w, viewData)
		}
	})

	fmt.Printf("Starting server at port 80\n")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
