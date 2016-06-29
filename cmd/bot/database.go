package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	
	log "github.com/Sirupsen/logrus"
)

var (
	database *sql.DB
)

func dbOpen(dbFile string) {
	var err error

	database, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}

	schemaCheckStmt := `
	SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='discord_stats';
	`
	discordStatsRows, err := database.Query(schemaCheckStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, schemaCheckStmt)
		panic(err)
	}
	if (dbRowCount(discordStatsRows) == 0) {
		log.Printf("Empty database, Creating Schema")
		dbInit()
	} else {
		log.Printf("Database found")
	}
	discordStatsRows.Close()
}

func dbInit() {	
	sqlStmt := `
	CREATE TABLE discord_stats (
		id integer not null primary key,
		username text not null,
		game text not null,
		time integer not null
	);
	`
	_, err := database.Exec(sqlStmt)
	if err != nil {
		log.Fatal("%q: %s\n", err, sqlStmt)
		panic(err)
	}
}

func dbIncGameEntry(username string, game string, increment int) {
	if len(username) <= 0 || len(game) <= 0 {
		return
	}

	log.Printf("Incrementing value for %s on %s by %d", username, game, increment);

	gameCheckStmt := "SELECT COUNT(*) FROM discord_stats WHERE username = ? AND game = ?;"
	gameRowCount, err := database.Query(gameCheckStmt, username, game)
	if err != nil {
		log.Fatal("%q: %s\n", err, gameCheckStmt)
		return
	}
	if (dbRowCount(gameRowCount) == 0) {
		log.Printf("Creating a row")
		_, err = database.Exec("INSERT INTO discord_stats(username, game, time) VALUES(?, ?, ?)", username, game, increment)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	var time int
	var id int

	err = database.QueryRow("SELECT id, time FROM discord_stats WHERE username = ? AND game = ?", username, game).Scan(&id, &time)
	if err != nil {
		log.Fatal("%q\n", err)
	}

	log.Printf("Got ID: %d, Current: %d, New: %d", id, time, time + increment)

	_, err = database.Exec("UPDATE discord_stats SET time = ? WHERE id = ?", time + increment, id)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func dbGetGameStats(username string) ([]string, []int){
	gameStatsStmt := "SELECT game, time FROM discord_stats WHERE username = ?;"
	rows, err := database.Query(gameStatsStmt, username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var games []string
	var times []int

	for rows.Next() {
		var game string
		var time int
		err = rows.Scan(&game, &time)
		if err != nil {
			log.Fatal(err)
		}
		games = append(games, game)
		times = append(times, time)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	return games, times
}

func dbRowCount(rows *sql.Rows) (count int) {
 	for rows.Next() {
    	err:= rows.Scan(&count)
    	if err != nil {
			panic(err)
		}
    }   
    return count
}
