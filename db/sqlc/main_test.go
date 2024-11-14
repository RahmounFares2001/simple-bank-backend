package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/RahmounFares2001/simple-bank-backend/util"
	_ "github.com/lib/pq"
)

// func newTestServer(t *testing.T, store db.Store) *Server {
// 	config := util.Config{
// 		TokenSymmetricKey:  util.RandomString(32),
// 		AccesTokenDuration: time.Minute,
// 	}
// 	server, err := NewServer(config, store)
// 	require.NoError(t, err)

// 	return server
// }

var testQueries *Queries

func TestMain(m *testing.M) {
	// env var
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db")
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
