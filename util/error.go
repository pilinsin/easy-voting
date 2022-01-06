package util

import (
	"errors"
	"fmt"
	"log"
)

func NewError(a ...interface{}) error {
	return errors.New(fmt.Sprintln(a...))
}
func AddError(err error, str string) error {
	return errors.New(fmt.Sprintln(str, ":", err))
}
func CheckError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
func RaiseError(a ...interface{}) {
	err := errors.New(fmt.Sprintln(a...))
	log.Panic(err)
}
func PrintError(err error) {
	if err != nil {
		log.Println(err)
	}
}
