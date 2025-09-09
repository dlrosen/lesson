package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type InstructorAvail struct {
	Instructor_ID  int
	Time_Period_ID int
	Seq_Nbr        int
	School_ID      int
	Day_Of_Week    string
	Start_Time     int
	End_Time       int
}

type ReturnInstructorAvail struct {
	ReturnInstructorAvail []InstructorAvail
	ReturnErr             error
	ReturnMsg             string
}

type ReturnDelInstructorAvail struct {
	ReturnAffected int64
	ReturnErr      error
	ReturnMsg      string
}

type ReturnMaxInstructorAvail struct {
	ReturnMaxSeq int
	ReturnErr    error
	ReturnMsg    string
}

func selectForInstructorTimeperiod(db *sql.DB, instructorID int, timePeriodID int) []InstructorAvail {
	rows, err := db.Query("SELECT instructor_id, time_period_id, seq_nbr, school_id, day_of_week, start_time, end_time FROM instructor_availability WHERE instructor_id = ? AND time_period_id = ?", instructorID, timePeriodID)

	studentCheckErr(err)

	defer rows.Close()

	instructorsAvail := make([]InstructorAvail, 0)

	for rows.Next() {
		ourInstructorAvail := InstructorAvail{}
		err = rows.Scan(&ourInstructorAvail.Instructor_ID, &ourInstructorAvail.Time_Period_ID, &ourInstructorAvail.Seq_Nbr, &ourInstructorAvail.School_ID, &ourInstructorAvail.Day_Of_Week, &ourInstructorAvail.Start_Time, &ourInstructorAvail.End_Time)
		if err != nil {
			log.Fatal(err)
		}

		instructorsAvail = append(instructorsAvail, ourInstructorAvail)
	}

	err = rows.Err()
	studentCheckErr(err)

	return instructorsAvail
}

func insertInstructorAvail(db *sql.DB, newinstructoravail InstructorAvail) []ReturnInstructorAvail {

	fmt.Printf("Starting insertInstructorAvail\n")
	returnStructure := make([]ReturnInstructorAvail, 0)
	instructorAvails := make([]InstructorAvail, 0)

	defaultReturnMsg := fmt.Sprintf("instructor_availability insert error for instructor_id:%d time_period_id:%d seq_nbr:%d", newinstructoravail.Instructor_ID, newinstructoravail.Time_Period_ID, newinstructoravail.Seq_Nbr)
	myReturnMsg := ""

	rows, err := db.Query("INSERT INTO instructor_availability (instructor_id, time_period_id, seq_nbr, school_id, day_of_week, start_time, end_time) VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING *", newinstructoravail.Instructor_ID, newinstructoravail.Time_Period_ID, newinstructoravail.Seq_Nbr, newinstructoravail.School_ID, newinstructoravail.Day_Of_Week, newinstructoravail.Start_Time, newinstructoravail.End_Time)

	if err != nil {
		myReturnMsg = fmt.Sprintf("db.Query %s", defaultReturnMsg)
	} else {
		defer rows.Close()

		for rows.Next() {
			ourInstructorAvail := InstructorAvail{}
			err = rows.Scan(&ourInstructorAvail.Instructor_ID, &ourInstructorAvail.Time_Period_ID, &ourInstructorAvail.Seq_Nbr, &ourInstructorAvail.School_ID, &ourInstructorAvail.Day_Of_Week, &ourInstructorAvail.Start_Time, &ourInstructorAvail.End_Time)

			if err != nil {
				myReturnMsg = fmt.Sprintf("rows.Scan %s", defaultReturnMsg)
			} else {
				instructorAvails = append(instructorAvails, ourInstructorAvail)
			}
		}
		err = rows.Err()

		if err != nil {
			myReturnMsg = fmt.Sprintf("row.Err %s", defaultReturnMsg)
		}
	}

	d := struct {
		ReturnInstructorAvail []InstructorAvail
		ReturnErr             error
		ReturnMsg             string
	}{
		ReturnInstructorAvail: instructorAvails,
		ReturnErr:             err,
		ReturnMsg:             myReturnMsg,
	}

	returnStructure = append(returnStructure, d)

	return returnStructure
}

func deleteInstructorAvail(db *sql.DB, delInstructorID int, delTimePeriodID int, delSeqNbr int) []ReturnDelInstructorAvail {

	returnStructure := make([]ReturnDelInstructorAvail, 0)
	defaultReturnMsg := fmt.Sprintf("instructor_availability delete error for instructor_id:%d time_period_id:%d seq_nbr:%d", delInstructorID, delTimePeriodID, delSeqNbr)
	myReturnMsg := ""
	var myAffected int64 = 0

	stmt, err := db.Prepare("DELETE FROM instructor_availability where instructor_id = ? AND time_period_id = ? AND seq_nbr = ?")

	if err != nil {
		myReturnMsg = fmt.Sprintf("db.Prepare %s", defaultReturnMsg)
	} else {
		defer stmt.Close()
		res, err := stmt.Exec(delInstructorID, delTimePeriodID, delSeqNbr)

		if err != nil {
			myReturnMsg = fmt.Sprintf("stmt.Exec %s", defaultReturnMsg)
		} else {
			affected, err := res.RowsAffected()

			if err != nil {
				myReturnMsg = fmt.Sprintf("res.RowsAffected %s", defaultReturnMsg)
			} else {
				myAffected = affected
			}

		}
	}

	d := struct {
		ReturnAffected int64
		ReturnErr      error
		ReturnMsg      string
	}{
		ReturnAffected: myAffected,
		ReturnErr:      err,
		ReturnMsg:      myReturnMsg,
	}

	returnStructure = append(returnStructure, d)

	return returnStructure
}

func selectMaxSeqInstructorAvail(db *sql.DB, instructorID int, timePeriodID int) []ReturnMaxInstructorAvail {

	returnStructure := make([]ReturnMaxInstructorAvail, 0)
	myMaxCnt := 0
	defaultReturnMsg := fmt.Sprintf("instructor_availability select max error for instructor_id:%d time_period_id:%d", instructorID, timePeriodID)
	myReturnMsg := ""

	rows, err := db.Query("SELECT IFNULL(MAX(seq_nbr),0) FROM instructor_availability WHERE instructor_id = ? AND time_period_id = ?", instructorID, timePeriodID)
	if err != nil {
		myReturnMsg = fmt.Sprintf("db.Query %s", defaultReturnMsg)
	} else {
		defer rows.Close()

		for rows.Next() {
			fmt.Printf("inside rows.Next \n")
			err = rows.Scan(&myMaxCnt)
			if err != nil {
				myReturnMsg = fmt.Sprintf("rows.Scan %s", defaultReturnMsg)
			}
		}
		err = rows.Err()

		if err != nil {
			myReturnMsg = fmt.Sprintf("row.Err %s", defaultReturnMsg)
		}
	}

	myMaxCnt = myMaxCnt + 1

	d := struct {
		ReturnMaxSeq int
		ReturnErr    error
		ReturnMsg    string
	}{
		ReturnMaxSeq: myMaxCnt,
		ReturnErr:    err,
		ReturnMsg:    myReturnMsg,
	}

	returnStructure = append(returnStructure, d)

	return returnStructure
}
