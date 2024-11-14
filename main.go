package main

import (
	"database/sql"
	"log"

	"github.com/RahmounFares2001/simple-bank-backend/api"
	db "github.com/RahmounFares2001/simple-bank-backend/db/sqlc"
	"github.com/RahmounFares2001/simple-bank-backend/util"
	_ "github.com/lib/pq"
)

func main() {
	// env var, passilou l path te3 l app.env file
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config", err)
	}

	// db cnx
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db")
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
