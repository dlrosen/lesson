package main

type PageFormat struct {
	PrevDisplay bool
	NextDisplay bool
	NewLow      int
}

func getPrevNext(button string, low int, maxcnt int, rows int) PageFormat {

	newlow := 0

	if button == "Next" {
		newlow = low + rows
	}

	if button == "Prev" {
		newlow = low - rows
	}

	prev := true
	next := true

	if newlow <= 0 {
		prev = false
	}

	if (maxcnt - (newlow + rows)) <= 0 {
		next = false
	}

	zPageFormat := PageFormat{
		PrevDisplay: prev,
		NextDisplay: next,
		NewLow:      newlow,
	}
	return zPageFormat
}

func getSearchFor(button string, searchFor string, searchForPrevNext string) string {

	newSearchFor := searchFor

	if button == "Next" {
		newSearchFor = searchForPrevNext
	}

	if button == "Prev" {
		newSearchFor = searchForPrevNext
	}

	return newSearchFor
}
