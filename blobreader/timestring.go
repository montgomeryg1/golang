package main

import (
	"fmt"
	"time"
)

func timeString(chSrch chan string, errorLog *log.logger, stopTimestring chan bool) {
	//t := time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC)
	t := time.Now().UTC().Add(-10 * time.Minute)
loop:
	for {
		y := t.Year()
		mon := t.Month()
		d := t.Day()
		h := t.Hour()
		m := t.Minute()
		tString := fmt.Sprintf("%d/%02d/%02d/%02d/%02d", y, mon, d, h, m)
		for {
			if t.Before(time.Now().Add(-1 * time.Minute)) {
				break
			} else {
				time.Sleep(1 * time.Minute)
			}
		}
		select {
		case <-stopTimestring:
			errorLog.Println("Stopping TimeString...")
			close(stopTimestring)
			break loop

		case <-time.After(time.Second):
		}
		chSrch <- tString
		t = t.Add(1 * time.Minute)
	}
}
