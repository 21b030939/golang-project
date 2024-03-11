package main

import (
	"database/sql"
	"flag"
	"github.com/21b030939/golang-project/pkg/schedule/model"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type config struct {
	Port string
	Env  string
	DB   struct {
		DSN string
	}
}

type application struct {
	config config
	models model.Models
}

func main() {
	var cfg config
	flag.StringVar(&cfg.Port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.DB.DSN, "db-dsn", "postgresql://postgres:password@localhost:5433/schedule?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
	}

	app.run()
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (app *application) run() {
	r := mux.NewRouter()

	v1 := r.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/schedules", app.createScheduleHandler).Methods("POST")
	v1.HandleFunc("/schedules/{scheduleId:[0-9]+}", app.getScheduleHandler).Methods("GET")
	v1.HandleFunc("/schedules/{scheduleId:[0-9]+}", app.updateScheduleHandler).Methods("PUT")
	v1.HandleFunc("/schedules/{scheduleId:[0-9]+}", app.deleteScheduleHandler).Methods("DELETE")

	log.Printf("Starting server on %s\n", app.config.Port)
	err := http.ListenAndServe(app.config.Port, r)
	log.Fatal(err)
}
