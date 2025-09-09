package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type TimePeriod struct {
	ID          int
	Description string
	StartDate   string
	EndDate     string
}

func insertTimePeriod(db *sql.DB, newtimeperiod TimePeriod) []TimePeriod {

	rows, err := db.Query("INSERT INTO time_period (time_period_id, description, start_date, end_date) VALUES (?, ?, ?, ?) RETURNING time_period_id, description, strftime('%Y-%m-%d', start_date), strftime('%Y-%m-%d', end_date)", nil, newtimeperiod.Description, newtimeperiod.StartDate, newtimeperiod.EndDate)
	timePeriodCheckErr(err)

	defer rows.Close()

	timeperiods := make([]TimePeriod, 0)

	for rows.Next() {
		ourTimePeriod := TimePeriod{}
		err = rows.Scan(&ourTimePeriod.ID, &ourTimePeriod.Description, &ourTimePeriod.StartDate, &ourTimePeriod.EndDate)
		if err != nil {
			log.Fatal(err)
		}

		timeperiods = append(timeperiods, ourTimePeriod)
	}

	err = rows.Err()
	timePeriodCheckErr(err)

	return timeperiods
}

func countForTimePeriod(db *sql.DB, searchString string) int {
	rows, err := db.Query("SELECT count(time_period_id) FROM time_period WHERE description like ?", searchString)
	timePeriodCheckErr(err)

	defer rows.Close()

	outMaxCnt := 0
	for rows.Next() {
		err = rows.Scan(&outMaxCnt)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rows.Err()
	timePeriodCheckErr(err)

	return outMaxCnt
}

func searchForTimePeriod(db *sql.DB, searchString string, limitPosition int, limitRows int) []TimePeriod {
	rows, err := db.Query("SELECT time_period_id, description, strftime('%Y-%m-%d', start_date), strftime('%Y-%m-%d', end_date) FROM time_period WHERE description like ? ORDER BY time_period_id LIMIT ?,?", searchString, limitPosition, limitRows)
	timePeriodCheckErr(err)

	defer rows.Close()

	timeperiods := make([]TimePeriod, 0)

	for rows.Next() {
		ourTimePeriod := TimePeriod{}
		err = rows.Scan(&ourTimePeriod.ID, &ourTimePeriod.Description, &ourTimePeriod.StartDate, &ourTimePeriod.EndDate)
		if err != nil {
			log.Fatal(err)
		}

		timeperiods = append(timeperiods, ourTimePeriod)
	}

	err = rows.Err()
	timePeriodCheckErr(err)

	return timeperiods
}

func searchForTimePeriodByID(db *sql.DB, searchID int) []TimePeriod {
	rows, err := db.Query("SELECT time_period_id, description, strftime('%Y-%m-%d', start_date), strftime('%Y-%m-%d', end_date) FROM time_period WHERE time_period_id = ? ", searchID)
	timePeriodCheckErr(err)

	defer rows.Close()

	timeperiods := make([]TimePeriod, 0)

	for rows.Next() {
		ourTimePeriod := TimePeriod{}
		err = rows.Scan(&ourTimePeriod.ID, &ourTimePeriod.Description, &ourTimePeriod.StartDate, &ourTimePeriod.EndDate)
		timePeriodCheckErr(err)

		timeperiods = append(timeperiods, ourTimePeriod)
	}

	err = rows.Err()
	timePeriodCheckErr(err)

	return timeperiods
}

func getTimePeriodListBox(db *sql.DB, timePeriodID int) []MyListBox {
	rows, err := db.Query("select case when t2.time_period_id is NULL then false else true end, t1.description||'-'||t1.time_period_id from time_period t1 LEFT OUTER JOIN time_period t2 ON t1.time_period_id = t2.time_period_id AND t2.time_period_id = ? ORDER BY t1.time_period_id", timePeriodID)

	timePeriodCheckErr(err)

	defer rows.Close()

	timePeriodListBox := make([]MyListBox, 0)

	for rows.Next() {
		ourTimePeriodListBox := MyListBox{}
		err = rows.Scan(&ourTimePeriodListBox.Select, &ourTimePeriodListBox.Option)
		if err != nil {
			log.Fatal(err)
		}

		timePeriodListBox = append(timePeriodListBox, ourTimePeriodListBox)
	}

	err = rows.Err()
	instructorCheckErr(err)

	return timePeriodListBox
}

func updateTimePeriod(db *sql.DB, ourTimePeriod TimePeriod) int64 {

	stmt, err := db.Prepare("UPDATE time_period set description = ?, start_date = ?, end_date = ? where time_period_id = ?")
	timePeriodCheckErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(ourTimePeriod.Description, ourTimePeriod.StartDate, ourTimePeriod.EndDate, ourTimePeriod.ID)
	timePeriodCheckErr(err)

	affected, err := res.RowsAffected()
	timePeriodCheckErr(err)

	return affected
}

func deleteTimePeriod(db *sql.DB, idToDelete int) int64 {

	stmt, err := db.Prepare("DELETE FROM time_period where time_period_id = ?")
	timePeriodCheckErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(idToDelete)
	timePeriodCheckErr(err)

	affected, err := res.RowsAffected()
	timePeriodCheckErr(err)

	return affected
}
