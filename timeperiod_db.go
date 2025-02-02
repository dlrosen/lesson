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
