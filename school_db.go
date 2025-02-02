package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type School struct {
	ID          int
	School_Name string
	Active      bool
}

type SchoolListBox struct {
	Select bool
	Option string
}

func insertSchool(db *sql.DB, newschool School) []School {

	rows, err := db.Query("INSERT INTO school (school_id, school_name, active) VALUES (?, ?, ?) RETURNING *", nil, newschool.School_Name, newschool.Active)
	schoolCheckErr(err)

	defer rows.Close()

	schools := make([]School, 0)

	for rows.Next() {
		ourSchool := School{}
		err = rows.Scan(&ourSchool.ID, &ourSchool.School_Name, &ourSchool.Active)
		if err != nil {
			log.Fatal(err)
		}

		schools = append(schools, ourSchool)
	}

	err = rows.Err()
	schoolCheckErr(err)

	return schools
}

func countForschool(db *sql.DB, searchString string) int {
	rows, err := db.Query("SELECT count(school_id) FROM school WHERE school_name like ?", searchString)
	schoolCheckErr(err)

	defer rows.Close()

	outMaxCnt := 0
	for rows.Next() {
		err = rows.Scan(&outMaxCnt)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rows.Err()
	schoolCheckErr(err)

	return outMaxCnt
}

func searchForschool(db *sql.DB, searchString string, limitPosition int, limitRows int) []School {
	rows, err := db.Query("SELECT school_id, school_name, active FROM school WHERE school_name like ? ORDER BY school_id LIMIT ?,?", searchString, limitPosition, limitRows)
	schoolCheckErr(err)

	defer rows.Close()

	schools := make([]School, 0)

	for rows.Next() {
		ourSchool := School{}
		err = rows.Scan(&ourSchool.ID, &ourSchool.School_Name, &ourSchool.Active)
		if err != nil {
			log.Fatal(err)
		}

		schools = append(schools, ourSchool)
	}

	err = rows.Err()
	schoolCheckErr(err)

	return schools
}

func searchForschoolByID(db *sql.DB, searchID int) []School {
	rows, err := db.Query("SELECT school_id, school_name, active FROM school WHERE school_id = ? ", searchID)
	schoolCheckErr(err)

	defer rows.Close()

	schools := make([]School, 0)

	for rows.Next() {
		ourSchool := School{}
		err = rows.Scan(&ourSchool.ID, &ourSchool.School_Name, &ourSchool.Active)
		schoolCheckErr(err)

		schools = append(schools, ourSchool)
	}

	err = rows.Err()
	schoolCheckErr(err)

	return schools
}

func getSchoolListBox(db *sql.DB, schoolID int) []SchoolListBox {
	rows, err := db.Query("select case when s2.school_id is NULL then false else true end, s1.school_name||'-'||s1.school_id from school s1 LEFT OUTER JOIN school s2 ON s1.school_id = s2.school_id AND s2.school_id = ? ORDER BY s1.school_id", schoolID)

	schoolCheckErr(err)

	defer rows.Close()

	schoolsListBox := make([]SchoolListBox, 0)

	for rows.Next() {
		ourSchoolListBox := SchoolListBox{}
		err = rows.Scan(&ourSchoolListBox.Select, &ourSchoolListBox.Option)
		if err != nil {
			log.Fatal(err)
		}

		schoolsListBox = append(schoolsListBox, ourSchoolListBox)
	}

	err = rows.Err()
	schoolCheckErr(err)

	return schoolsListBox
}

func updateSchool(db *sql.DB, ourSchool School) int64 {

	stmt, err := db.Prepare("UPDATE school set school_name = ?, active = ? where school_id = ?")
	schoolCheckErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(ourSchool.School_Name, ourSchool.Active, ourSchool.ID)
	schoolCheckErr(err)

	affected, err := res.RowsAffected()
	schoolCheckErr(err)

	return affected
}

func deleteSchool(db *sql.DB, idToDelete int) int64 {

	stmt, err := db.Prepare("DELETE FROM school where school_id = ?")
	schoolCheckErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(idToDelete)
	schoolCheckErr(err)

	affected, err := res.RowsAffected()
	schoolCheckErr(err)

	return affected
}
