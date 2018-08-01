package errorHandle

import (
	"errors"
	"fmt"
	"log"
)

func CatchLoadDataError(err *error) {
	if r := recover(); r != nil {
		log.Println("scanner data to struct error: ", r)
		*err = errors.New(fmt.Sprint(r))
		panic(r)
	}
}

func CatchError() {
	if r := recover(); r != nil {
		log.Println("catch error :", r)
		panic(r)
	}
}
