package app

import "time"

// nowDate returns today's date as YYYY-MM-DD.
func nowDate() string {
	return time.Now().Format("2006-01-02")
}
