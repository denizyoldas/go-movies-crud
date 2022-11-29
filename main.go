package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"strconv"
)

type Movie struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Year int `json:"year"`
	Director *Director `json:"director"`
	Actors []string `json:"actors"`
}

type Director struct {
	ID int `json:"id"`
	Name string `json:"name"`
}

var movies = []Movie{
	{ID: 1, Title: "The Shawshank Redemption", Year: 1994, Director: &Director{ID: 1, Name: "Frank Darabont"}, Actors: []string{"Tim Robbins", "Morgan Freeman", "Bob Gunton"}},
	{ID: 2, Title: "The Godfather", Year: 1972, Director: &Director{ID: 2, Name: "Francis Ford Coppola"}, Actors: []string{"Marlon Brando", "Al Pacino", "James Caan"}},
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	http.ListenAndServe(":8080", r)
	fmt.Println("Server started on port 8080")
}


func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		id, _ := strconv.Atoi(params["id"])
		if item.ID == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Movie{})
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = len(movies) + 1
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		id, _ := strconv.Atoi(params["id"])
		if item.ID == id {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = id
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		id, _ := strconv.Atoi(params["id"])
		if item.ID == id {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}