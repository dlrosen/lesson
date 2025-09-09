package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type InstructorAvailRows struct {
	Row_Nbr        int
	Instructor_ID  int
	Time_Period_ID int
	Seq_Nbr        int
	Schools        []MyListBox
	Days_Of_Week   []MyListBox
	Start_Hour     []MyListBox
	Start_Minute   []MyListBox
	End_Hour       []MyListBox
	End_Minute     []MyListBox
}

func instructorAvailSearchTemplate(
	w http.ResponseWriter,
	r *http.Request,
	db *sql.DB,
) {
	zInstructureListBox := getInstructorListBox(db, 1)
	zTimePeriodListBox := getTimePeriodListBox(db, 1)

	d := struct {
		InstructorList []MyListBox
		TimePeriodList []MyListBox
	}{
		InstructorList: zInstructureListBox,
		TimePeriodList: zTimePeriodListBox,
	}

	tpl.ExecuteTemplate(w, "instructor_avail_search.gohtml", d)
}

func searchInstructorAvail(
	w http.ResponseWriter,
	r *http.Request,
	db *sql.DB,
) {
	myInstructorID := splitForID(r.FormValue("form_instructorlist"))
	myInstructorValue := splitForValue(r.FormValue("form_instructorlist"))

	myTimePeriodID := splitForID(r.FormValue("form_timeperiodlist"))
	myTimePeriodValue := splitForValue(r.FormValue("form_timeperiodlist"))

	listInstructorAvail(myInstructorValue, myInstructorID, myTimePeriodValue, myTimePeriodID, w, db)
}

func updateInstructorAvail(
	w http.ResponseWriter,
	r *http.Request,
	db *sql.DB,
) {

	myMaxCnt := convStringToNbr(r.FormValue("form_max_rows"))
	myInstructorID := convStringToNbr(r.FormValue("form_instructor_id"))
	myInstructorValue := r.FormValue("form_instructor_value")
	myTimePeriodID := convStringToNbr(r.FormValue("form_time_period_id"))
	myTimePeriodValue := r.FormValue("form_time_period_value")

	cnt := 1
	for {
		cnt_string := strconv.Itoa(cnt)
		myDelete := r.FormValue(fmt.Sprintf("form_delete-%s", cnt_string))
		mySeqNbr := convStringToNbr(r.FormValue(fmt.Sprintf("form_seq-%s", cnt_string)))
		myDay := r.FormValue(fmt.Sprintf("form_daylist-%s", cnt_string))
		mySchoolID := splitForID(r.FormValue(fmt.Sprintf("form_schoollist-%s", cnt_string)))

		myStartHour := convStringToNbr(r.FormValue(fmt.Sprintf("form_starthourlist-%s", cnt_string)))
		myStartMinute := convStringToNbr(r.FormValue(fmt.Sprintf("form_startminutelist-%s", cnt_string)))
		myStartTime := (myStartHour * 100) + myStartMinute

		myEndHour := convStringToNbr(r.FormValue(fmt.Sprintf("form_endhourlist-%s", cnt_string)))
		myEndMinute := convStringToNbr(r.FormValue(fmt.Sprintf("form_endminutelist-%s", cnt_string)))
		myEndTime := (myEndHour * 100) + myEndMinute

		if myDelete == "yes" {
			fmt.Printf("DELETE path\n")
		} else {
			if mySeqNbr > 0 {
				fmt.Printf("UPDATE path\n")
			} else {
				fmt.Printf("INSERT path\n")

				// add logic to get MAX Seq_Nbr
				myMaxSeqNbr := 0
				returnMaxInstructorAvail := selectMaxSeqInstructorAvail(db, myInstructorID, myTimePeriodID)
				for _, myReturnMaxInstructorAvail := range returnMaxInstructorAvail {
					fmt.Printf("Return %s \n", myReturnMaxInstructorAvail.ReturnMsg)
					fmt.Printf("Return Max Seq %d \n", myReturnMaxInstructorAvail.ReturnMaxSeq)
					myMaxSeqNbr = myReturnMaxInstructorAvail.ReturnMaxSeq
				}

				newInstructorAvail := InstructorAvail{
					Instructor_ID:  myInstructorID,
					Time_Period_ID: myTimePeriodID,
					Seq_Nbr:        myMaxSeqNbr,
					School_ID:      mySchoolID,
					Day_Of_Week:    myDay,
					Start_Time:     myStartTime,
					End_Time:       myEndTime,
				}

				cacheMutex.Lock()
				returnInstructorAvail := insertInstructorAvail(db, newInstructorAvail)
				cacheMutex.Unlock()
				for _, myReturnInstructorAvail := range returnInstructorAvail {
					fmt.Printf("Return %s \n", myReturnInstructorAvail.ReturnMsg)
				}
			}
		}
		//		if lookup_delete == "yes" {
		//         if seqnbr > 0
		//         {
		//            call the delete function
		//         } else {
		//            skipping delete
		//         }
		//     	} else {
		//         if seqnbr > 0 {
		//           check for differences from DB
		//           if differences then UPDATE
		//         } else {
		//           check for differences from defaults
		//           if differences then INSERT
		//         }
		//		}

		if cnt >= myMaxCnt {
			break
		}
		cnt = cnt + 1
	}

	listInstructorAvail(myInstructorValue, myInstructorID, myTimePeriodValue, myTimePeriodID, w, db)
}

func listInstructorAvail(
	myInstructorValue string,
	myInstructorID int,
	myTimePeriodValue string,
	myTimePeriodID int,
	w http.ResponseWriter,
	db *sql.DB,
) {
	zPage := make([]InstructorAvailRows, 0)

	myStartHour := 0
	myEndHour := 0
	myStartMinute := 0
	myEndMinute := 0

	myRowNbr := 0

	cacheMutex.Lock()
	instructorTimeperiod := selectForInstructorTimeperiod(db, myInstructorID, myTimePeriodID)
	cacheMutex.Unlock()

	for _, myInstructorTimeperiod := range instructorTimeperiod {
		myRowNbr = myRowNbr + 1
		fullPage := InstructorAvailRows{}

		myStartHour = myInstructorTimeperiod.Start_Time / 100
		myEndHour = myInstructorTimeperiod.End_Time / 100

		myStartMinute = myInstructorTimeperiod.Start_Time - (myStartHour * 100)
		myEndMinute = myInstructorTimeperiod.End_Time - (myEndHour * 100)

		fullPage.Row_Nbr = myRowNbr
		fullPage.Instructor_ID = myInstructorTimeperiod.Instructor_ID
		fullPage.Time_Period_ID = myInstructorTimeperiod.Time_Period_ID
		fullPage.Seq_Nbr = myInstructorTimeperiod.Seq_Nbr
		fullPage.Schools = getSchoolListBox(db, myInstructorTimeperiod.School_ID)
		fullPage.Days_Of_Week = getDaysListBox(myInstructorTimeperiod.Day_Of_Week)
		fullPage.Start_Hour = getTimeListBox(myStartHour, 23)
		fullPage.Start_Minute = getTimeListBox(myStartMinute, 59)
		fullPage.End_Hour = getTimeListBox(myEndHour, 23)
		fullPage.End_Minute = getTimeListBox(myEndMinute, 59)

		zPage = append(zPage, fullPage)
	}

	cnt := 0
	for {
		myRowNbr = myRowNbr + 1
		emptyPage := InstructorAvailRows{}

		emptyPage.Row_Nbr = myRowNbr
		emptyPage.Instructor_ID = myInstructorID
		emptyPage.Time_Period_ID = myTimePeriodID
		emptyPage.Seq_Nbr = 0
		emptyPage.Schools = getSchoolListBox(db, 0)
		emptyPage.Days_Of_Week = getDaysListBox("Mon")
		emptyPage.Start_Hour = getTimeListBox(0, 23)
		emptyPage.Start_Minute = getTimeListBox(0, 59)
		emptyPage.End_Hour = getTimeListBox(0, 23)
		emptyPage.End_Minute = getTimeListBox(0, 59)

		zPage = append(zPage, emptyPage)

		if cnt >= 2 {
			break
		}
		cnt = cnt + 1
	}

	yMessage := fmt.Sprintf("Searched for %s %s", myInstructorValue, myTimePeriodValue)

	d := struct {
		InstructorAvailList []InstructorAvailRows
		YourMessage         string
		MaxRows             int
		InstructorID        int
		InstructorValue     string
		TimePeriodID        int
		TimePeriodValue     string
	}{
		InstructorAvailList: zPage,
		YourMessage:         yMessage,
		MaxRows:             myRowNbr,
		InstructorID:        myInstructorID,
		InstructorValue:     myInstructorValue,
		TimePeriodID:        myTimePeriodID,
		TimePeriodValue:     myTimePeriodValue,
	}

	tpl.ExecuteTemplate(w, "instructor_avail_list.gohtml", d)
}

func convStringToNbr(theString string) int {
	myNbr, err := strconv.Atoi(theString)
	if err != nil {
		myNbr = 0
	}

	return myNbr
}

func splitForID(theString string) int {
	myList := strings.Split(theString, "-")
	myNbr, err := strconv.Atoi(myList[1])

	if err != nil {
		myNbr = 0
	}

	return myNbr
}

func splitForValue(theString string) string {
	myList := strings.Split(theString, "-")
	myValue := myList[0]

	return myValue
}

//func instructorAvailCheckErr(err error, w http.ResponseWriter, message string) {
//
//	if err != nil {
//		fmt.Printf("message: %s \n", message)
//		fmt.Printf("err.Error(): %s \n", err.Error())
//
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//}
