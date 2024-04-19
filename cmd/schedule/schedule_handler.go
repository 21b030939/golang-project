package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/21b030939/golang-project/pkg/schedule/model"
	"github.com/21b030939/golang-project/pkg/schedule/validator"
)

func (app *application) createScheduleHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Discipline string `json:"discipline"`
		Cabinet    string `json:"cabinet"`
		TimePeriod int    `json:"timePeriod"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		log.Println(err)
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	schedule := &model.Schedule{
		Discipline: input.Discipline,
		Cabinet:    input.Cabinet,
		TimePeriod: input.TimePeriod,
	}

	err = app.models.Schedules.Insert(schedule)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"schedule": schedule}, nil)
}

func (app *application) getScheduleList(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Discipline          string
		TimePeriodValueFrom int
		TimePeriodValueTo   int
		model.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	// Use our helpers to extract the title and nutrition value range query string values, falling back to the
	// defaults of an empty string and an empty slice, respectively, if they are not provided
	// by the client.
	input.Discipline = app.readStrings(qs, "discipline", "")
	input.TimePeriodValueFrom = app.readInt(qs, "timePeriodFrom", 0, v)
	input.TimePeriodValueTo = app.readInt(qs, "timePeriodTo", 0, v)

	// Ge the page and page_size query string value as integers. Notice that we set the default
	// page value to 1 and default page_size to 20, and that we pass the validator instance
	// as the final argument.
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	// Extract the sort query string value, falling back to "id" if it is not provided
	// by the client (which will imply an ascending sort on menu ID).
	input.Filters.Sort = app.readStrings(qs, "sort", "id")

	// Add the supported sort value for this endpoint to the sort safelist.
	// name of the column in the database.
	input.Filters.SortSafeList = []string{
		// ascending sort values
		"id", "discipline", "time_period",
		// descending sort values
		"-id", "-discipline", "-time_period",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	schedules, metadata, err := app.models.Schedules.GetAll(input.Discipline, input.TimePeriodValueFrom, input.TimePeriodValueTo, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"menus": schedules, "metadata": metadata}, nil)
}

func (app *application) getScheduleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	schedule, err := app.models.Schedules.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"schedule": schedule}, nil)
}

func (app *application) updateScheduleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	schedule, err := app.models.Schedules.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Discipline *string `json:"discipline"`
		Cabinet    *string `json:"cabinet"`
		TimePeriod *int    `json:"timePeriod"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
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

	v := validator.New()

	if model.ValidateSchedule(v, schedule); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Schedules.Update(schedule)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"schedule": schedule}, nil)
}

func (app *application) deleteScheduleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Schedules.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "success"}, nil)
}