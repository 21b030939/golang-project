package model

import (
	// "database/sql"
	"errors"
	"log"
	"github.com/jmoiron/sqlx"
)

type Discipline struct {
	Id          string `json:"id"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Credits     string `json:"credits"`
}

type DisciplineModel struct {
	DB       *sqlx.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

var disciplines = []Discipline{
	{
		Id:          "1",
		Name:        "Discrete Structures",
		Description: "This course covers introductory topics in discrete mathematics such as sets, mathematical reasoning and proofs,combinatorial counting methods and generating functions, basics of number theory, and basics of graph theory.",
		Credits:     "3",
	},
	{
		Id:          "2",
		Name:        "Calculus II",
		Description: "This course is the second part of a mathematics course. It contains the following chapters: antiderivatives; definite integrals; applications of definite integrals; differentiable calculus of functions of two or more variables; multiple integrals.",
		Credits:     "3",
	},
	{
		Id:          "3",
		Name:        "Programming Principles I",
		Description: "C++ was designed with systems programming and embedded, resource-constrained software and large systems in mind, with performance, efficiency, and flexibility of use as its design highlights.",
		Credits:     "4",
	},
	{
		Id:          "4",
		Name:        "Android Development",
		Description: "Tools and APIs required building applications for the Android platform using the Android SDK. User interface designs for mobile devices and unique user interactions using multi-touch technologies. Object-oriented design using model-view-controller paradigm, memory management, Java (Kotlin) programming language. Other topics include: object-oriented database API, animation, multi-threading and performance considerations.",
		Credits:     "3",
	},
	{
		Id:          "5",
		Name:        "Golang Application Development",
		Description: "Go is a statically typed, compiled high-level programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. It is syntactically similar to C, but also has memory safety, garbage collection, structural typing, and CSP-style concurrency.",
		Credits:     "4",
	},
}

func GetDisciplines() []Discipline {
	return disciplines
}

func GetDiscipline(id string) (*Discipline, error) {
	for _, r := range disciplines {
		if r.Id == id {
			return &r, nil
		}
	}
	return nil, errors.New("discipline not found")
}
