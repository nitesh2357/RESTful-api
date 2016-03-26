package main

import (
	"github.com/drone/routes"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
   // "fmt"
   // "time"
    
	
)

type Request struct {
	Email string `json:"email"`
	Zip string `json:"zip"`
	Country string `json:"country"`
	Profession string	`json:"profession"`
	Favorite_color string `json:"favorite_color"`
	Is_smoking string	`json:"is_smoking"`
	Favorite_sport string `json:"favorite_sport"`
	Food Foods `json:"food"`
	Music Musics `json:"music"`
	Movie Movies `json:"movie"`
	Travel Travels `json:"travel"`
}

type Foods struct {
	Type string `json:"type"`
	Drink_alcohol string `json:"drink_alcohol"`
} 

type Musics struct {
	Spotify_user_id string `json:"spotify_user_id"`
}

type Movies struct{
	Tv_shows []string `json:"tv_shows"`
	Movies []string `json:"movies"`
}

type Travels struct{
	Flight Flights `json:"flight"`
}

type Flights struct{
	Seat string `json:"seat"`
}

var post []Request

func PostProfile(w http.ResponseWriter, r *http.Request){
	
	var req Request
	decoder := json.NewDecoder(r.Body)
	error := decoder.Decode(&req)
	if error != nil{
		http.Error(w,"Json Error",http.StatusBadRequest)
		return
	}
	post = append(post,req)
    //fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", 201, "Created")
    //fmt.Fprintf(w, "Date: %s", time.Now().UTC().Format(http.TimeFormat))
	w.WriteHeader(http.StatusCreated)
     
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	email := params.Get(":email")
	index := isPresent(email)
	profile, err := json.Marshal(post[index])
	if err != nil {
		w.WriteHeader(505)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(profile))
    

}

func isPresent(email string) int{

	for i , value := range post {
		if value.Email == email{
			return i
		}
	}
	return -1
}



func DelProfile(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()
	email := params.Get(":email")
	index := isPresent(email)
	post = append(post[:index],post[index+1:]...)
	w.WriteHeader(204)
}

func PutProfile(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()
	email := params.Get(":email")
	index := isPresent(email)
	if index == -1 {
				w.WriteHeader(204)
				return
	}
	var new_data map[string]interface{}
	updateJson, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(updateJson    , &new_data)

	

	if err != nil{
		http.Error(w,"Json Error",http.StatusBadRequest)
		return
	}

	new_profile,_ := json.Marshal(post[index])

	var new_profile1 map[string]interface{}
	_ = json.Unmarshal(new_profile, &new_profile1)

	new_profile1 = updateprofile(new_data, new_profile1)

	new_profile,_ = json.Marshal(new_profile1)
	_ = json.Unmarshal(new_profile, &post[index])
		

	w.WriteHeader(204)
}

func updateprofile(updateDate map[string]interface{}, profileData map[string]interface{}) map[string]interface{}{
	for i := range updateDate {
		profileData[i] = updateDate[i]
	}
	return profileData
}

func main() {
	mux := routes.New()

	mux.Get("/profile/:email", GetProfile)
	mux.Post("/profile", PostProfile)
	mux.Del("/profile/:email",DelProfile)
	mux.Put("/profile/:email",PutProfile)
	http.Handle("/", mux)
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}


