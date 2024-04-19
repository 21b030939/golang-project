package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/21b030939/golang-project/pkg/schedule/validator"
)

type Schedule struct {
	Id         string `json:"id"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
	Discipline string `json:"discipline"`
	Cabinet    string `json:"cabinet"`
	TimePeriod int    `json:"timePeriod"`
}

type ScheduleModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m ScheduleModel) GetAll(description string, from, to int, filters Filters) ([]*Schedule, Metadata, error) {

	// Retrieve all menu items from the database.
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, created_at, updated_at, discipline, cabinet, time_period
		FROM schedule
		WHERE (LOWER(discipline) = LOWER($1) OR $1 = '')
		AND (time_period >= $2 OR $2 = 0)
		AND (time_period <= $3 OR $3 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5
		`,
		filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Organize our four placeholder parameter values in a slice.
	args := []interface{}{description, from, to, filters.limit(), filters.offset()}

	// log.Println(query, description, from, to, filters.limit(), filters.offset())
	// Use QueryContext to execute the query. This returns a sql.Rows result set containing
	// the result.
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	// Importantly, defer a call to rows.Close() to ensure that the result set is closed
	// before GetAll returns.
	defer func() {
		if err := rows.Close(); err != nil {
			m.ErrorLog.Println(err)
		}
	}()

	// Declare a totalRecords variable
	totalRecords := 0

	var schedules []*Schedule
	for rows.Next() {
		var schedule Schedule
		err := rows.Scan(&totalRecords, &schedule.Id, &schedule.CreatedAt, &schedule.UpdatedAt, &schedule.Discipline, &schedule.Cabinet, &schedule.TimePeriod)
		if err != nil {
			return nil, Metadata{}, err
		}

		// Add the Movie struct to the slice
		schedules = append(schedules, &schedule)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Generate a Metadata struct, passing in the total record count and pagination parameters
	// from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// If everything went OK, then return the slice of the movies and metadata.
	return schedules, metadata, nil
}

func (m ScheduleModel) Insert(schedule *Schedule) error {
	// Insert a new schedule item into the database.
	query := `
		INSERT INTO schedule (discipline, cabinet, time_period) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{schedule.Discipline, schedule.Cabinet, schedule.TimePeriod}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&schedule.Id, &schedule.CreatedAt, &schedule.UpdatedAt)
}

func (m ScheduleModel) Get(id int) (*Schedule, error) {
	// Return an error if the ID is less than 1.
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Retrieve a specific schedule item based on its ID.
	query := `
		SELECT id, created_at, updated_at, discipline, cabinet, time_period
		FROM schedule
		WHERE id = $1
		`
	var schedule Schedule
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&schedule.Id, &schedule.CreatedAt, &schedule.UpdatedAt, &schedule.Discipline, &schedule.Cabinet, &schedule.TimePeriod)
	if err != nil {
		return nil, fmt.Errorf("cannot retrive menu with id: %v, %w", id, err)
	}
	return &schedule, nil
}

func (m ScheduleModel) Update(schedule *Schedule) error {
	// Update a specific schedule item in the database.
	query := `
		UPDATE schedule
		SET discipline = $1, cabinet = $2, time_period = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND updated_at = $5
		RETURNING updated_at
	`
	args := []interface{}{schedule.Discipline, schedule.Cabinet, schedule.TimePeriod, schedule.Id, schedule.UpdatedAt}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&schedule.UpdatedAt)
}

func (m ScheduleModel) Delete(id int) error {
	// Return an error if the ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}

	// Delete a specific schedule item from the database.
	query := `
		DELETE FROM schedule
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func ValidateSchedule(v *validator.Validator, schedule *Schedule) {
	// Check if the discipline field is empty.
	v.Check(schedule.Discipline != "", "discipline", "must be provided")
	// Check if the discipline field is not more than 100 characters.
	v.Check(len(schedule.Discipline) <= 100, "discipline", "must not be more than 100 bytes long")
	// Check if the description field is not more than 1000 characters.
	v.Check(len(schedule.Cabinet) <= 1000, "description", "must not be more than 1000 bytes long")
	// Check if the nutrition value is not more than 10000.
	v.Check(schedule.TimePeriod <= 6, "nutritionValue", "must not be more than 10000")
}
