package utils

import (
	"errors"
	"fmt"
)

func ErrorPrintf(str string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(str, args...))
}

func ErrorPrint(str string) error {
	return errors.New(str)
}
