package lib

import (
	"errors"
	"time"
)

func SetupTimeLocale(location string) error {
	loc, err := time.LoadLocation(location)

	if err != nil {
		return errors.New("Unable to set time locale to: " + location + "\n" + err.Error())
	}

	time.Local = loc

	return err
}
