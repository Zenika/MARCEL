package media

import (
	"net/http"
	//"github.com/gorilla/mux"
	//"strconv"
	"io/ioutil"
	"encoding/json"
	"errors"
	"log"
)

var Medias []Media
var y interface{}

func LoadMedias() {
	content, err := ioutil.ReadFile("data/media.config.json")
	check(err)

	/*err = json.Unmarshal(content, &Medias)
	check(err)

	for _, m := range Medias {
		for _, p:=range m.Plugins {
			err = json.Unmarshal(content, &p.PropValues.X)
			check(err)
		}
	}*/

	json.Unmarshal([]byte(content), &y)
	log.Print(y)

	log.Print("Medias config is loaded...")
	log.Print(Medias)
}

func GetMedia(idMedia int) (*Media, error) {
	for _, m := range Medias {
		if idMedia == m.ID {
			return &m, nil
		}
	}

	return nil, errors.New("NO_MEDIA_FOUND")
}

func HandleGetMedia(w http.ResponseWriter, r *http.Request) {
	/*vars := mux.Vars(r)
	f := vars["idMedia"]
	idMedia, _ := strconv.Atoi(f)

	m, err := GetMedia(idMedia)
	if err != nil {
		writeResponseWithError(w, http.StatusNotFound)
		return
	}*/

	//b, err := json.Marshal(*m)
	b, err := json.Marshal(y)
	if err != nil {
		writeResponseWithError(w, http.StatusNotFound)
		return
	}

	w.Write([]byte(b))

}

func writeResponseWithError(w http.ResponseWriter, errorCode int) {
	w.WriteHeader(errorCode)
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
		panic(e)
	}
}
