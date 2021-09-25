package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	sentry "github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
)

func init() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://18e7bd4ed00d4ebda6bee172125d3118@o1010291.ingest.sentry.io/5974739",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	defer sentry.Flush(2 * time.Second)

	sentry.CaptureMessage("It works!")

	var filename string = "general.log"
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		sentry.CaptureMessage(err.Error())
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
	Image   string
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Println("root method: ", r.Method)
	if r.Method == "GET" {
		http.Redirect(w, r, "/home", http.StatusFound)
	} else {
		sentry.CaptureMessage("Method received different than GET")
		log.WithFields(log.Fields{
			"uri": "*",
		}).Error("Method received different than GET")
		http.Error(w, "Only GET methods are supported", http.StatusNotFound)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home method: ", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("src/html/home.html")
		t.Execute(w, nil)
	} else {
		sentry.CaptureMessage("Method received different than GET")
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

					im := generateRandomImage()
					respon1 := Respon{
						Boolean: b,
						Name:    name,
						Realm:   realm,
						Region:  region,
						Level:   level,
						Image:   im,
					}

					t := template.Must(template.ParseFiles("src/html/response.html"))
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
		sentry.CaptureMessage("Method received different than GET")
		log.WithFields(log.Fields{
			"uri": "/response",
		}).Error("Method received different than GET")
		http.Error(w, "Only GET methods are supported", http.StatusNotFound)
	}
}

func characterError(w http.ResponseWriter, r *http.Request) {
	fmt.Println("index method: ", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("src/html/error.html")
		t.Execute(w, nil)
	} else {
		sentry.CaptureMessage("Method received different than GET")
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

	fs := http.FileServer(http.Dir("src/html/css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))
	fsi := http.FileServer(http.Dir("src/html/images"))
	http.Handle("/images/", http.StripPrefix("/images/", fsi))

	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"Port": os.Getenv("PORT"),
		}).Fatal("Error on ListenAndServe")
		fmt.Println("Fatal on listenAndServe")
	}
}

func generateRandomImage() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	var number int
	for ok := true; ok; ok = true {
		number = r1.Intn(5)
		if number != 0 {
			break
		}
	}

	str := "images/" + strconv.Itoa(number) + ".gif"
	fmt.Println(str)
	return str
}
