package main

import (
	"encoding/json"
	"github.com/21b030939/golang-project/pkg/schedule/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *application) createScheduleHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Discipline        string `json:"discipline"`
		Cabinet           string `json:"cabinet"`
		TimePeriod        string `json:"timePeriod"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	schedule := &model.Schedule{
		Discipline:    input.Discipline,
		Cabinet:       input.Cabinet,
		TimePeriod:    input.TimePeriod,
	}

	err = app.models.Schedules.Insert(schedule)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, schedule)
}

func (app *application) getScheduleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["scheduleId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	menu, err := app.models.Schedules.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, menu)
}

func (app *application) updateScheduleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["scheduleId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	schedule, err := app.models.Schedules.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Discipline *string `json:"discipline"`
		Cabinet    *string `json:"cabinet"`
		TimePeriod *string `json:"timePeriod"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Discipline != nil {
		schedule.Discipline = *input.Discipline
	}

	if input.Cabinet != nil {
		schedule.Cabinet = *input.Cabinet
	}

	if input.TimePeriod != nil {
		schedule.TimePeriod = *input.TimePeriod
	}

	err = app.models.Schedules.Update(schedule)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, schedule)
}

func (app *application) deleteScheduleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["scheduleId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid menu ID")
		return
	}

	err = app.models.Schedules.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}