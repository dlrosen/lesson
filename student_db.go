package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Student struct {
	ID          int
	First_Name  string
	Last_Name   string
	Email       string
	School_ID   int
	Active      bool
	School_Name string
}

func insertStudent(db *sql.DB, newstudent Student, schoolName string) []Student {

	rows, err := db.Query("INSERT INTO student (student_id, first_name, last_name, email, school_id, active) VALUES (?, ?, ?, ?, ?, ?) RETURNING *", nil, newstudent.First_Name, newstudent.Last_Name, newstudent.Email, newstudent.School_ID, newstudent.Active)
	studentCheckErr(err)

	defer rows.Close()

	students := make([]Student, 0)

	for rows.Next() {
		ourStudent := Student{}
		err = rows.Scan(&ourStudent.ID, &ourStudent.First_Name, &ourStudent.Last_Name, &ourStudent.Email, &ourStudent.School_ID, &ourStudent.Active)
		if err != nil {
			log.Fatal(err)
		}

		ourStudent.School_Name = schoolName
		students = append(students, ourStudent)
	}

	err = rows.Err()
	studentCheckErr(err)

	return students
}

func countForstudent(db *sql.DB, searchString string) int {
	rows, err := db.Query("SELECT count(s1.student_id) FROM student s1, school s2 WHERE s1.school_id = s2.school_id AND (s1.first_name like ? OR s1.last_name like ?) ", searchString, searchString)
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
	rows, err := db.Query("SELECT s1.student_id, s1.first_name, s1.last_name, s1.email, s1.school_id, s1.active, s2.school_name FROM student s1, school s2 WHERE s1.school_id = s2.school_id AND (s1.first_name like ? OR s1.last_name like ?) ORDER BY student_id LIMIT ?,?", searchString, searchString, limitPosition, limitRows)
	studentCheckErr(err)

	defer rows.Close()

	students := make([]Student, 0)

	for rows.Next() {
		ourStudent := Student{}
		err = rows.Scan(&ourStudent.ID, &ourStudent.First_Name, &ourStudent.Last_Name, &ourStudent.Email, &ourStudent.School_ID, &ourStudent.Active, &ourStudent.School_Name)
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
	rows, err := db.Query("SELECT s1.student_id, s1.first_name, s1.last_name, s1.email, s1.school_id, s1.active, s2.school_name FROM student s1, school s2 WHERE s1.school_id = s2.school_id AND s1.student_id = ? ", searchID)
	studentCheckErr(err)

	defer rows.Close()

	students := make([]Student, 0)

	for rows.Next() {
		ourStudent := Student{}
		err = rows.Scan(&ourStudent.ID, &ourStudent.First_Name, &ourStudent.Last_Name, &ourStudent.Email, &ourStudent.School_ID, &ourStudent.Active, &ourStudent.School_Name)
		studentCheckErr(err)

		students = append(students, ourStudent)
	}

	err = rows.Err()
	studentCheckErr(err)

	return students
}

func updateStudent(db *sql.DB, ourStudent Student) int64 {

	stmt, err := db.Prepare("UPDATE student set first_name = ?, last_name = ?, email = ?, school_id = ?, active = ? where student_id = ?")
	studentCheckErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(ourStudent.First_Name, ourStudent.Last_Name, ourStudent.Email, ourStudent.School_ID, ourStudent.Active, ourStudent.ID)
	studentCheckErr(err)

	affected, err := res.RowsAffected()
	studentCheckErr(err)

	return affected
}

func deleteStudent(db *sql.DB, idToDelete int) int64 {

	stmt, err := db.Prepare("DELETE FROM student where student_id = ?")
	studentCheckErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(idToDelete)
	studentCheckErr(err)

	affected, err := res.RowsAffected()
	studentCheckErr(err)

	return affected
}
