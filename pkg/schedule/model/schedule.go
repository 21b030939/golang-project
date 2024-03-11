package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Schedule struct {
	Id         string `json:"id"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
	Discipline string `json:"discipline"`
	Cabinet    string `json:"cabinet"`
	TimePeriod string `json:"timePeriod"`
}

type ScheduleModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m ScheduleModel) Insert(schedule *Schedule) error {
	// Insert a new menu item into the database.
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
	// Retrieve a specific menu item based on its ID.
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
		return nil, err
	}
	return &schedule, nil
}

func (m ScheduleModel) Update(schedule *Schedule) error {
	// Update a specific menu item in the database.
	query := `
		UPDATE schedule
		SET discipline = $1, cabinet = $2, time_period = $3
		WHERE id = $4
		RETURNING updated_at
		`
	args := []interface{}{schedule.Discipline, schedule.Cabinet, schedule.TimePeriod, schedule.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&schedule.UpdatedAt)
}

func (m ScheduleModel) Delete(id int) error {
	// Delete a specific menu item from the database.
	query := `
		DELETE FROM schedule
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
