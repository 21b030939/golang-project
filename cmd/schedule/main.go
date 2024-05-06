package main

import (
	// "database/sql"
	"flag"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/21b030939/golang-project/pkg/jsonlog"
	"github.com/21b030939/golang-project/pkg/schedule/model"
	"github.com/21b030939/golang-project/pkg/schedule/model/filter"
	"github.com/21b030939/golang-project/pkg/vcs"
	_ "github.com/lib/pq"
)

// Set version of application corresponding to value of vcs.Version.
var (
	version = vcs.Version()
)

type config struct {
	Port int
	Env  string
	Fill bool
	DB   struct {
		DSN string
	}
}

type application struct {
	config config
	models model.Models
	logger *jsonlog.Logger
	wg     sync.WaitGroup
}

func main() {
	var cfg config
	flag.BoolVar(&cfg.Fill, "fill", false, "Fill Database with some data")
	flag.IntVar(&cfg.Port, "port", 8080, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.DB.DSN, "db-dsn", "host=db port=5432 user=postgres dbname=schedule password=postgres sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	// Init logger
	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintError(err, nil)
		return
	}
	// Defer a call to db.Close() so that the connection pool is closed before the main()
	// function exits.
	defer func() {
		if err := db.Close(); err != nil {
			logger.PrintFatal(err, nil)
		}
	}()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
		logger: logger,
	}

	if cfg.Fill{
		err = filler.PopulateDatabase(app.models)
		if err != nil{
			logger.PrintFatal(err, nil)
			return
		}
	}

	if err := app.serve(); err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sqlx.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config // struct.
	db, err := sqlx.Open("postgres", cfg.DB.DSN)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
