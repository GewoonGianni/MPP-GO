package main

// De hidden line is "Your mission, should you choose to accept it," uit Mission Impossible :)

import (
	"flag"
	"fmt"
	"log"
	"os"

	"database/sql"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "movies.db")
	// ----- Ook mogelijk om zonder GCC te doen, hierbij kan je "sqlite3" vervangen door "sqlite" in de statemant hierboven, "github.com/mattn/go-sqlite3" vervangen door "modernc.org/sqlite" bij import en het command "go mod tidy" in de terminal te typen. Dit is "pure go" maar gebruikt in de go.mod alsnog de mattn import onder require. (mits ik dit goed begreep)
	checkError(err)

	arguments := os.Args[1:] // The first element is the path to the command, so we can skip that

	addCommand := flag.NewFlagSet("add", flag.ExitOnError)
	addImdbId := addCommand.String("imdbid", "tt0000001", "The IMDb ID of a movie or series")
	addTitle := addCommand.String("title", "Carmencita", "The movie's or series' title")
	addYear := addCommand.Int("year", 1894, "The movie's or series' year of release")
	addImdbRating := addCommand.Float64("rating", 5.7, "The movie's or series' rating on IMDb")

	detailsCommand := flag.NewFlagSet("details", flag.ExitOnError)
	detailsImdbId := detailsCommand.String("imdbid", "tt0000001", "The IMDb ID of a movie or series")

	deleteCommand := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteImdbId := deleteCommand.String("imdbid", "tt0000001", "The IMDb ID of a movie or series")

	switch arguments[0] {
	case "add":
		addCommand.Parse(arguments[1:])

		stmt, err := db.Prepare("INSERT INTO movies(IMDb_id, Title, Rating, Year) values(?,?,?,?)")
		checkError(err)

		stmt.Exec(*addImdbId, *addTitle, *addImdbRating, *addYear)

		fmt.Println("IMDb id:", *addImdbId)
		fmt.Println("Title:", *addTitle)
		fmt.Println("Rating:", *addImdbRating)
		fmt.Println("Year:", *addYear)
	case "list":
		rows, err := db.Query("SELECT Title FROM movies")
		checkError(err)
		var Title string

		for rows.Next() {
			err = rows.Scan(&Title)
			checkError(err)
			fmt.Println(Title)
		}

		rows.Close()
	case "details":
		detailsCommand.Parse(arguments[1:])

		query := fmt.Sprintf("SELECT * FROM movies WHERE IMDb_id = '%s'", *detailsImdbId)

		rows, err := db.Query(query)
		checkError(err)
		var IMDb_id string
		var Title string
		var Rating float64
		var year int

		for rows.Next() {
			err = rows.Scan(&IMDb_id, &Title, &Rating, &year)
			checkError(err)
			fmt.Println("IMDb id:", IMDb_id)
			fmt.Println("Title:", Title)
			fmt.Println("Rating:", Rating)
			fmt.Println("Year:", year)
		}

		rows.Close()
	case "delete":
		deleteCommand.Parse(arguments[1:])

		stmt, err := db.Prepare("delete from movies where IMDb_id=?")
		checkError(err)

		stmt.Exec(*deleteImdbId)

		fmt.Println("Movie deleted")
	}
	db.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
