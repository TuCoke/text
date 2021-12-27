package util

import (
	"log"
	"time"
)

func Retry(attempts int, sleep time.Duration, f func(interface{}) error, i interface{}) error {
	if err := f(i); err != nil {
		if s, ok := err.(stop); ok {
			return s.error
		}
		if attempts--; attempts > 0 {
			log.Printf("retry func error: %s. attemps #%d after %s\n", err.Error(), attempts, sleep)
			time.Sleep(sleep)
			return Retry(attempts, 2*sleep, f,i)
		}
		return err
	}
	return nil
}

type stop struct {
	error
}

func NoRetryError(err error) stop {
	return stop{err}
}