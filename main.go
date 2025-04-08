package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

type Subject struct {
	ID      int    `db:"id" json:"id"`
	Subject string `db:"subject" json:"subject"`
	SubjectUUID string `db:"subject_uuid" json:"subject_uuid"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

func main() {
	e := echo.New()

	db, err := sqlx.Open("sqlite3", "./database.sqlite")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	e.GET("/subject", func(c echo.Context) error {
		subjectName := c.QueryParam("slug")
		query := fmt.Sprintf("SELECT * FROM subjects WHERE subject_uuid='%s'", subjectName)
		log.Println("Executing query:", query)

		var subjects []Subject
		err := db.Select(&subjects, query)
		if err != nil {
			log.Println("SQL error:", err)
			return c.String(http.StatusInternalServerError, "Database error")
		}
		return c.JSON(http.StatusOK, subjects)
	})

	e.GET("/subject-safe", func(c echo.Context) error {
		subjectName := c.QueryParam("slug")
		query := `SELECT * FROM subjects WHERE subject_uuid = ?`
		log.Println("[Safe] Executing parameterized query")

		var subjects []Subject
		err := db.Select(&subjects, query, subjectName)
		if err != nil {
			log.Println("SQL error:", err)
			return c.String(http.StatusInternalServerError, "Database error")
		}
		return c.JSON(http.StatusOK, subjects)
	})

	e.Logger.Fatal(e.Start(":8000"))
}
