package main

import (
	"fmt"
)

type MyListBox struct {
	Select bool
	Option string
}

func getDaysListBox(pick string) []MyListBox {

	myDayListBox := make([]MyListBox, 0)
	days := []string{"Sun", "Mon", "Tue", "Wed", "Thr", "Fri", "Sat"}

	for _, myDay := range days {
		ourDayListBox := MyListBox{}
		if pick == myDay {
			ourDayListBox.Select = true
		} else {
			ourDayListBox.Select = false
		}

		ourDayListBox.Option = myDay
		myDayListBox = append(myDayListBox, ourDayListBox)
	}

	return myDayListBox
}

func getTimeListBox(pick int, limit int) []MyListBox {

	theListBox := make([]MyListBox, 0)

	cnt := 0
	for {
		ourListBox := MyListBox{}
		if pick == cnt {
			ourListBox.Select = true
		} else {
			ourListBox.Select = false
		}
		ourListBox.Option = fmt.Sprintf("%02d", cnt)
		theListBox = append(theListBox, ourListBox)
		if cnt >= limit {
			break
		}
		cnt = cnt + 1
	}

	return theListBox
}
