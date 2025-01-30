package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type School struct {
	ID          int
	School_Name string
}

func insertSchool(db *sql.DB, newschool School) []School {

	rows, err := db.Query("INSERT INTO school (id, school_name) VALUES (?, ?) RETURNING *", nil, newschool.School_Name)
	schoolCheckErr(err)

	defer rows.Close()

	schools := make([]School, 0)

	for rows.Next() {
		ourSchool := School{}
		err = rows.Scan(&ourSchool.ID, &ourSchool.School_Name)
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
	rows, err := db.Query("SELECT count(id) FROM school WHERE school_name like ?", searchString)
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
	rows, err := db.Query("SELECT id, school_name FROM school WHERE school_name like ? ORDER BY id LIMIT ?,?", searchString, limitPosition, limitRows)
	schoolCheckErr(err)

	defer rows.Close()

	schools := make([]School, 0)

	for rows.Next() {
		ourSchool := School{}
		err = rows.Scan(&ourSchool.ID, &ourSchool.School_Name)
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
	rows, err := db.Query("SELECT id, school_name FROM school WHERE id = ? ", searchID)
	schoolCheckErr(err)

	defer rows.Close()

	schools := make([]School, 0)

	for rows.Next() {
		ourSchool := School{}
		err = rows.Scan(&ourSchool.ID, &ourSchool.School_Name)
		schoolCheckErr(err)

		schools = append(schools, ourSchool)
	}

	err = rows.Err()
	schoolCheckErr(err)

	return schools
}

func updateSchool(db *sql.DB, ourSchool School) int64 {

	stmt, err := db.Prepare("UPDATE school set school_name = ? where id = ?")
	schoolCheckErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(ourSchool.School_Name, ourSchool.ID)
	schoolCheckErr(err)

	affected, err := res.RowsAffected()
	schoolCheckErr(err)

	return affected
}

func deleteSchool(db *sql.DB, idToDelete int) int64 {

	stmt, err := db.Prepare("DELETE FROM school where id = ?")
	schoolCheckErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(idToDelete)
	schoolCheckErr(err)

	affected, err := res.RowsAffected()
	schoolCheckErr(err)

	return affected
}
