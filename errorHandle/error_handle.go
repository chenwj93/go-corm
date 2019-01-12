package errorHandle

import (
	"log"
	"runtime/debug"
)

func CatchLoadDataError(err *error) {
	if r := recover(); r != nil {
		log.Println("scanner data to struct error: ", r)
		debug.PrintStack()
	}
}

func CatchError() {
	if r := recover(); r != nil {
		log.Println("catch error :", r)
		debug.PrintStack()
	}
}
