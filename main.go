package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "movies"
)

type Movie struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Year     int       `json:"year"`
	Director *Director `json:"director"`
}

type Director struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

var db *sql.DB

func main() {
	r := mux.NewRouter()

	sqlOpen, dbErr := sql.Open("postgres", psqlInfo)
	db = sqlOpen

	if dbErr != nil {
		fmt.Println(dbErr)
	}

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Server started on port 3232")

	err := http.ListenAndServe(":3232", r)
	if err != nil {
		panic("Server Error: " + err.Error())
	}
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query(`
		SELECT m.id, m.title, m.year, d.id AS director_id, d.name AS director_name
		FROM movies m
		INNER JOIN directors d ON m.director_id = d.id
	`)

	if err != nil {
		// Handle the error by returning a JSON response
		json.NewEncoder(w).Encode(map[string]bool{"success": false})
		return
	}

	defer rows.Close()

	var mvList []Movie

	for rows.Next() {
		var m Movie
		// Initialize the Director field
		m.Director = &Director{}
		// Modify the Scan to include the new director fields
		if err := rows.Scan(&m.ID, &m.Title, &m.Year, &m.Director.ID, &m.Director.Name); err != nil {
			// Handle the error by returning a JSON response
			fmt.Println(err)
			json.NewEncoder(w).Encode(map[string]bool{"success": false})
			return
		}

		// Append the movie to the list
		mvList = append(mvList, m)
	}

	// Encode the movie list as JSON and send it in the response
	json.NewEncoder(w).Encode(mvList)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	paramId := params["id"]

	var movie Movie
	movie.Director = &Director{}

	errQ := db.QueryRow(`
		SELECT m.id, m.title, m.year, d.id AS director_id, d.name AS director_name
		FROM movies m
		INNER JOIN directors d ON m.director_id = d.id WHERE m.id = $1`, paramId).Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Director.ID, &movie.Director.Name)

	if errQ != nil {
		fmt.Println(errQ)
		json.NewEncoder(w).Encode(map[string]bool{"success": false})
		return
	}

	json.NewEncoder(w).Encode(movie)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type CreatePayload struct {
		Title      string `json:"title"`
		Year       int    `json:"year"`
		DirectorId int    `json:"directorId"`
	}
	var movie CreatePayload
	_ = json.NewDecoder(r.Body).Decode(&movie)

	defer db.Close()

	exc, err := db.Exec("INSERT INTO movies (title, year, director_id) values ($1,$2,$3)", movie.Title, movie.Year, movie.DirectorId)

	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(map[string]bool{"success": false})
		return
	}

	json.NewEncoder(w).Encode(exc)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	type UpdatePayload struct {
		Title      string `json:"title"`
		Year       int    `json:"year"`
		DirectorId int    `json:"directorId"`
	}

	var movie UpdatePayload
	_ = json.NewDecoder(r.Body).Decode(&movie)

	exc, err := db.Exec("UPDATE movies SET title = $1, year = $2, director_id = $3 WHERE id = $4",
		movie.Title, movie.Year, movie.DirectorId, id)

	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(map[string]bool{"success": false})
		return
	}

	fmt.Println(exc)

	json.NewEncoder(w).Encode(exc)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	exc, err := db.Exec("DELETE FROM movies WHERE id = $1", id)

	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(map[string]bool{"success": false})
		return
	}

	json.NewEncoder(w).Encode(exc)
}
