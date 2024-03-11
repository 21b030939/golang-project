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
	var input model.Schedule
	if err := app.readJSON(r, &input); err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := app.models.Schedules.Insert(&input); err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, input)
}

func (app *application) getScheduleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["scheduleId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	schedule, err := app.models.Schedules.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Schedule not found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, schedule)
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
		app.respondWithError(w, http.StatusNotFound, "Schedule not found")
		return
	}

	var input struct {
		Discipline *string `json:"discipline"`
		Cabinet    *string `json:"cabinet"`
		TimePeriod *string `json:"timePeriod"`
	}

	if err := app.readJSON(r, &input); err != nil {
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

	if err := app.models.Schedules.Update(schedule); err != nil {
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
		app.respondWithError(w, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	err = app.models.Schedules.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) readJSON(r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}

	return nil
}
