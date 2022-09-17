package components

import "log"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func doLog(format string, v ...any) {
	log.Printf(format, v...)
}
