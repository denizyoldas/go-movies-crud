# Go Movies App

This is a simple app that uses the [Gorilla Mux](https://github.com/gorilla/mux) package to create a RESTful API. The API is used to manage a list of movies. The API is used by a simple web app that allows users to view the list of movies and add new movies to the list.

## Getting Started

To get started, you will need to install the following:

- [Go](https://golang.org/dl/)

## Running the App

To run the app, you will need to open two terminal windows. In the first terminal window, run the following command:

```bash
go run main.go
```

## Requesting the API

To request the API, you can use a tool like [Postman](https://www.getpostman.com/). The following are the endpoints that are available:

- GET /movies - Returns a list of movies
- GET /movies/{id} - Returns a single movie
- POST /movies - Adds a new movie
- PUT /movies/{id} - Updates a movie
- DELETE /movies/{id} - Deletes a movie

## Author

- **[Deniz Aksu](https://denizaksu.dev/)**
