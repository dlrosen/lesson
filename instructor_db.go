package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Instructor struct {
	ID              int
	Instructor_Name string
	Active          bool
}

func insertInstructor(db *sql.DB, newinstructor Instructor) []Instructor {

	rows, err := db.Query("INSERT INTO instructor (instructor_id, instructor_name, active) VALUES (?, ?, ?) RETURNING *", nil, newinstructor.Instructor_Name, newinstructor.Active)
	instructorCheckErr(err)

	defer rows.Close()

	instructors := make([]Instructor, 0)

	for rows.Next() {
		ourInstructor := Instructor{}
		err = rows.Scan(&ourInstructor.ID, &ourInstructor.Instructor_Name, &ourInstructor.Active)
		if err != nil {
			log.Fatal(err)
		}

		instructors = append(instructors, ourInstructor)
	}

	err = rows.Err()
	instructorCheckErr(err)

	return instructors
}

func countForinstructor(db *sql.DB, searchString string) int {
	rows, err := db.Query("SELECT count(instructor_id) FROM instructor WHERE instructor_name like ?", searchString)
	instructorCheckErr(err)

	defer rows.Close()

	outMaxCnt := 0
	for rows.Next() {
		err = rows.Scan(&outMaxCnt)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rows.Err()
	instructorCheckErr(err)

	return outMaxCnt
}

func searchForinstructor(db *sql.DB, searchString string, limitPosition int, limitRows int) []Instructor {
	rows, err := db.Query("SELECT instructor_id, instructor_name, active FROM instructor WHERE instructor_name like ? ORDER BY instructor_id LIMIT ?,?", searchString, limitPosition, limitRows)
	instructorCheckErr(err)

	defer rows.Close()

	instructors := make([]Instructor, 0)

	for rows.Next() {
		ourInstructor := Instructor{}
		err = rows.Scan(&ourInstructor.ID, &ourInstructor.Instructor_Name, &ourInstructor.Active)
		if err != nil {
			log.Fatal(err)
		}

		instructors = append(instructors, ourInstructor)
	}

	err = rows.Err()
	instructorCheckErr(err)

	return instructors
}

func searchForinstructorByID(db *sql.DB, searchID int) []Instructor {
	rows, err := db.Query("SELECT instructor_id, instructor_name, active FROM instructor WHERE instructor_id = ? ", searchID)
	instructorCheckErr(err)

	defer rows.Close()

	instructors := make([]Instructor, 0)

	for rows.Next() {
		ourInstructor := Instructor{}
		err = rows.Scan(&ourInstructor.ID, &ourInstructor.Instructor_Name, &ourInstructor.Active)
		instructorCheckErr(err)

		instructors = append(instructors, ourInstructor)
	}

	err = rows.Err()
	instructorCheckErr(err)

	return instructors
}

func getInstructorListBox(db *sql.DB, instructorID int) []MyListBox {
	rows, err := db.Query("select case when i2.instructor_id is NULL then false else true end, i1.instructor_name||'-'||i1.instructor_id from instructor i1 LEFT OUTER JOIN instructor i2 ON i1.instructor_id = i2.instructor_id AND i2.instructor_id = ? ORDER BY i1.instructor_id", instructorID)

	instructorCheckErr(err)

	defer rows.Close()

	instructorsListBox := make([]MyListBox, 0)

	for rows.Next() {
		ourInstructorListBox := MyListBox{}
		err = rows.Scan(&ourInstructorListBox.Select, &ourInstructorListBox.Option)
		if err != nil {
			log.Fatal(err)
		}

		instructorsListBox = append(instructorsListBox, ourInstructorListBox)
	}

	err = rows.Err()
	instructorCheckErr(err)

	return instructorsListBox
}

func updateInstructor(db *sql.DB, ourInstructor Instructor) int64 {

	stmt, err := db.Prepare("UPDATE instructor set instructor_name = ?, active = ? where instructor_id = ?")
	instructorCheckErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(ourInstructor.Instructor_Name, ourInstructor.Active, ourInstructor.ID)
	instructorCheckErr(err)

	affected, err := res.RowsAffected()
	instructorCheckErr(err)

	return affected
}

func deleteInstructor(db *sql.DB, idToDelete int) int64 {

	stmt, err := db.Prepare("DELETE FROM instructor where instructor_id = ?")
	instructorCheckErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(idToDelete)
	instructorCheckErr(err)

	affected, err := res.RowsAffected()
	instructorCheckErr(err)

	return affected
}
