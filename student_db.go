package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Student struct {
	ID         int
	First_Name string
	Last_Name  string
	Email      string
}

func insertStudent(db *sql.DB, newstudent Student) []Student {

	rows, err := db.Query("INSERT INTO student (id, first_name, last_name, email) VALUES (?, ?, ?, ?) RETURNING *", nil, newstudent.First_Name, newstudent.Last_Name, newstudent.Email)
	studentCheckErr(err)

	defer rows.Close()

	students := make([]Student, 0)

	for rows.Next() {
		ourStudent := Student{}
		err = rows.Scan(&ourStudent.ID, &ourStudent.First_Name, &ourStudent.Last_Name, &ourStudent.Email)
		if err != nil {
			log.Fatal(err)
		}

		students = append(students, ourStudent)
	}

	err = rows.Err()
	studentCheckErr(err)

	return students
}

func countForstudent(db *sql.DB, searchString string) int {
	rows, err := db.Query("SELECT count(id) FROM student WHERE first_name like ? OR last_name like ?", searchString, searchString)
	studentCheckErr(err)

	defer rows.Close()

	outMaxCnt := 0
	for rows.Next() {
		err = rows.Scan(&outMaxCnt)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rows.Err()
	studentCheckErr(err)

	return outMaxCnt
}

func searchForstudent(db *sql.DB, searchString string, limitPosition int, limitRows int) []Student {
	rows, err := db.Query("SELECT id, first_name, last_name, email FROM student WHERE first_name like ? OR last_name like ? ORDER BY id LIMIT ?,?", searchString, searchString, limitPosition, limitRows)
	studentCheckErr(err)

	defer rows.Close()

	students := make([]Student, 0)

	for rows.Next() {
		ourStudent := Student{}
		err = rows.Scan(&ourStudent.ID, &ourStudent.First_Name, &ourStudent.Last_Name, &ourStudent.Email)
		if err != nil {
			log.Fatal(err)
		}

		students = append(students, ourStudent)
	}

	err = rows.Err()
	studentCheckErr(err)

	return students
}

func searchForstudentByID(db *sql.DB, searchID int) []Student {
	rows, err := db.Query("SELECT id, first_name, last_name, email FROM student WHERE id = ? ", searchID)
	studentCheckErr(err)

	defer rows.Close()

	students := make([]Student, 0)

	for rows.Next() {
		ourStudent := Student{}
		err = rows.Scan(&ourStudent.ID, &ourStudent.First_Name, &ourStudent.Last_Name, &ourStudent.Email)
		studentCheckErr(err)

		students = append(students, ourStudent)
	}

	err = rows.Err()
	studentCheckErr(err)

	return students
}

func updateStudent(db *sql.DB, ourStudent Student) int64 {

	stmt, err := db.Prepare("UPDATE student set first_name = ?, last_name = ?, email = ? where id = ?")
	studentCheckErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(ourStudent.First_Name, ourStudent.Last_Name, ourStudent.Email, ourStudent.ID)
	studentCheckErr(err)

	affected, err := res.RowsAffected()
	studentCheckErr(err)

	return affected
}

func deleteStudent(db *sql.DB, idToDelete int) int64 {

	stmt, err := db.Prepare("DELETE FROM student where id = ?")
	studentCheckErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(idToDelete)
	studentCheckErr(err)

	affected, err := res.RowsAffected()
	studentCheckErr(err)

	return affected
}
