package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func timePeriodCreateTemplate(
	w http.ResponseWriter,
	r *http.Request,
) {
	tpl.ExecuteTemplate(w, "timeperiod_create.gohtml", nil)
}

func timePeriodSearchTemplate(
	w http.ResponseWriter,
	r *http.Request,
) {

	d := struct {
		YourMessage string
	}{
		YourMessage: "",
	}

	tpl.ExecuteTemplate(w, "timeperiod_search.gohtml", d)
}

func timePeriodModifyTemplate(
	w http.ResponseWriter,
	r *http.Request,
	db *sql.DB,
) {
	yMessage := ""
	pDescription := ""
	pStartDate := ""
	pEndDate := ""
	pID := 0

	id, err := strconv.Atoi(r.PathValue("id"))
	timePeriodCheckErr(err)
	timePeriods := searchForTimePeriodByID(db, id)

	for index, myTimePeriod := range timePeriods {
		if index == 0 {
			pDescription = myTimePeriod.Description
			pStartDate = myTimePeriod.StartDate
			pEndDate = myTimePeriod.EndDate
			pID = myTimePeriod.ID
		}
	}

	d := struct {
		YourMessage string
		ID          int
		Description string
		StartDate   string
		EndDate     string
	}{
		YourMessage: yMessage,
		ID:          pID,
		Description: pDescription,
		StartDate:   pStartDate,
		EndDate:     pEndDate,
	}

	tpl.ExecuteTemplate(w, "timeperiod_modify.gohtml", d)
}

func createTimePeriod(
	w http.ResponseWriter,
	r *http.Request,
	db *sql.DB,
) {
	newtimePeriod := TimePeriod{
		Description: r.FormValue("form_description"),
		StartDate:   r.FormValue("form_start_date"),
		EndDate:     r.FormValue("form_end_date"),
	}

	pID := 0

	cacheMutex.Lock()
	timePeriods := insertTimePeriod(db, newtimePeriod)
	cacheMutex.Unlock()

	for index, myTimePeriod := range timePeriods {
		if index == 0 {
			pID = myTimePeriod.ID
		}
	}

	yMessage := fmt.Sprintf("timePeriod %d Created", pID)

	d := struct {
		TimePeriodList []TimePeriod
		YourMessage    string
		Low            int
		PrevDisplay    bool
		NextDisplay    bool
		SearchFor      string
	}{
		TimePeriodList: timePeriods,
		YourMessage:    yMessage,
		Low:            0,
		PrevDisplay:    false,
		NextDisplay:    false,
		SearchFor:      "",
	}

	tpl.ExecuteTemplate(w, "timeperiod_list.gohtml", d)
}

func modifyTimePeriod(
	w http.ResponseWriter,
	r *http.Request,
	db *sql.DB,
) {

	fmt.Printf("Form Button = %s \n", r.FormValue("form_button"))

	id, err := strconv.Atoi(r.FormValue("form_id"))

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	yMessage := ""

	if r.FormValue("form_button") == "Delete" {
		cacheMutex.Lock()
		deleteTimePeriod(db, id)
		cacheMutex.Unlock()

		yMessage = fmt.Sprintf("timePeriod %d Deleted", id)
	} else {
		ztimePeriod := TimePeriod{
			ID:          id,
			Description: r.FormValue("form_description"),
			StartDate:   r.FormValue("form_start_date"),
			EndDate:     r.FormValue("form_end_date"),
		}

		cacheMutex.Lock()
		updateTimePeriod(db, ztimePeriod)
		cacheMutex.Unlock()

		yMessage = fmt.Sprintf("timePeriod %d Modified", id)
	}

	cacheMutex.Lock()
	timePeriods := searchForTimePeriodByID(db, id)
	cacheMutex.Unlock()

	d := struct {
		TimePeriodList []TimePeriod
		YourMessage    string
		Low            int
		PrevDisplay    bool
		NextDisplay    bool
		SearchFor      string
	}{
		TimePeriodList: timePeriods,
		YourMessage:    yMessage,
		Low:            0,
		PrevDisplay:    false,
		NextDisplay:    false,
		SearchFor:      "",
	}

	tpl.ExecuteTemplate(w, "timeperiod_list.gohtml", d)
}

func searchTimePeriod(
	w http.ResponseWriter,
	r *http.Request,
	db *sql.DB,
) {
	low := 0
	rows := 100
	maxCount := 0

	low, err := strconv.Atoi(r.FormValue("form_low"))

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	NewSearchFor := getSearchFor(r.FormValue("form_button"), r.FormValue("form_timeperiodsearch"), r.FormValue("form_prev_next_searchfor"))

	cacheMutex.Lock()
	maxCount = countForTimePeriod(db, NewSearchFor)
	cacheMutex.Unlock()

	zPage := getPrevNext(r.FormValue("form_button"), low, maxCount, rows)

	cacheMutex.Lock()
	timePeriods := searchForTimePeriod(db, NewSearchFor, zPage.NewLow, rows)
	cacheMutex.Unlock()

	yMessage := fmt.Sprintf("Searched for %s", NewSearchFor)

	d := struct {
		TimePeriodList []TimePeriod
		YourMessage    string
		Low            int
		PrevDisplay    bool
		NextDisplay    bool
		SearchFor      string
	}{
		TimePeriodList: timePeriods,
		YourMessage:    yMessage,
		Low:            zPage.NewLow,
		PrevDisplay:    zPage.PrevDisplay,
		NextDisplay:    zPage.NextDisplay,
		SearchFor:      NewSearchFor,
	}

	tpl.ExecuteTemplate(w, "timeperiod_list.gohtml", d)
}

func timePeriodCheckErr(err error) {
	if err != nil {
		//d := struct {
		//	YourMessage string
		//}{
		//	YourMessage: "Error",
		//}

		//tpl.ExecuteTemplate(w, "timeperiod_search.gohtml", d)

		log.Fatal(err)
		//fmt.Printf("Error: %s", err.Error())
		return
	}
}
