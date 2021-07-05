package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

var PORT string = "8080"

func init() {
	var filename string = "general.log"
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		log.Println(err)
	} else {
		log.SetOutput(f)
	}
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.TraceLevel)
}

type Respon struct {
	Boolean string
	Name    string
	Realm   string
	Region  string
	Level   int
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Println("root method: ", r.Method)
	if r.Method == "GET" {
		http.Redirect(w, r, "/home", http.StatusFound)
	} else {
		log.WithFields(log.Fields{
			"uri": "*",
		}).Error("Method received different than GET")
		http.Error(w, "Only GET methods are supported", http.StatusNotFound)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home method: ", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("html/home.html")
		t.Execute(w, nil)
	} else {
		log.WithFields(log.Fields{
			"uri": "/home",
		}).Error("Method received different than GET")
		http.Error(w, "Only GET methods are supported", http.StatusNotFound)
	}
}

func response(w http.ResponseWriter, r *http.Request) {
	fmt.Println("response method: ", r.Method)
	if r.Method == "GET" {
		name := r.FormValue("input_name")
		realm := r.FormValue("input_realm")
		region := r.FormValue("input_region")

		if name != "" && realm != "" && region != "" {
			if (len(name) >= 2 && len(name) <= 12) &&
				(len(realm) >= 2 && len(realm) <= 24) {

				bearerToken := auth()
				level, responseCode := isChar60(realm, name, region, bearerToken)

				if level > 0 && responseCode == 200 {

					fmt.Println(level + responseCode)

					b := "No"
					if level == 60 {
						b = "Yes"
					}

					respon1 := Respon{
						Boolean: b,
						Name:    name,
						Realm:   realm,
						Region:  region,
						Level:   level,
					}

					t := template.Must(template.ParseFiles("html/response.html"))
					t.Execute(w, respon1)

					log.WithFields(log.Fields{
						"character":  respon1.Name,
						"realm":      respon1.Realm,
						"region":     respon1.Region,
						"char_level": respon1.Level,
					}).Info("Method received different than GET")

				} else {
					log.WithFields(log.Fields{
						"uri":           "/response",
						"level":         level,
						"response_code": responseCode,
					}).Error("Level or response_code error")
					http.Redirect(w, r, "/character_error", http.StatusFound)
				}
			} else {
				log.WithFields(log.Fields{
					"uri":          "/response",
					"length_name":  len(name),
					"length_realm": len(realm),
				}).Error("Wrong length of fields")
				http.Redirect(w, r, "/character_error", http.StatusFound)
			}
		}
	} else {
		log.WithFields(log.Fields{
			"uri": "/response",
		}).Error("Method received different than GET")
		http.Error(w, "Only GET methods are supported", http.StatusNotFound)
	}
}

func characterError(w http.ResponseWriter, r *http.Request) {
	fmt.Println("index method: ", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("html/error.html")
		t.Execute(w, nil)
	} else {
		log.WithFields(log.Fields{
			"uri": "/character_error",
		}).Error("Method received different than GET")
		http.Error(w, "Only GET methods are supported", http.StatusNotFound)
	}
}

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/character_error", characterError)
	http.HandleFunc("/home", home)
	http.HandleFunc("/response", response)

	fs := http.FileServer(http.Dir("html/css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"Port": PORT,
		}).Fatal("Error on ListenAndServe")
		fmt.Println("Fatal on listenAndServe")
	}
}
