package filler

import (
	model "github.com/21b030939/golang-project/pkg/schedule/model"
)

func PopulateDatabase(models model.Models) error {
	for _, schedule := range schedules {
		models.Schedules.Insert(&schedule)
	}
	// TODO: Implement disciplines pupulation
	// TODO: Implement the relationship between disciplines and schedules
	return nil
}

var schedules = []model.Schedule{
	{Discipline: "Calculus", Cabinet: "256", TimePeriod: 3},
	{Discipline: "Discrete Structures", Cabinet: "269a", TimePeriod: 3},
	{Discipline: "Linear Algebra", Cabinet: "259", TimePeriod: 3},
	{Discipline: "Calculus II", Cabinet: "283", TimePeriod: 3},
	{Discipline: "Statistics", Cabinet: "351", TimePeriod: 4},
	{Discipline: "Programming Principles", Cabinet: "269", TimePeriod: 4},
	{Discipline: "OOP", Cabinet: "269", TimePeriod: 4},
	{Discipline: "Programming Principles II", Cabinet: "461", TimePeriod: 4},
	{Discipline: "Databases", Cabinet: "428", TimePeriod: 3},
	{Discipline: "Algorithms", Cabinet: "Konaev Hall", TimePeriod: 3},
	{Discipline: "Android Development", Cabinet: "272", TimePeriod: 3},
	{Discipline: "Advanced Android", Cabinet: "272", TimePeriod: 3},
	{Discipline: "Golang", Cabinet: "383", TimePeriod: 3},
	{Discipline: "Spring", Cabinet: "444", TimePeriod: 3},
	{Discipline: "Web Development", Cabinet: "359", TimePeriod: 4},
	{Discipline: "IOS Development", Cabinet: "283", TimePeriod: 3},
	{Discipline: "Software Development", Cabinet: "461", TimePeriod: 3},
}
