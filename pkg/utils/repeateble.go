package utils

import (
	"fmt"
	"time"
)

func DoWithTries(fn func() error, attemps int, delay time.Duration) error {
	const op = "utils.repeatable.DoWithTries"

	if attemps == 0 {
		return fmt.Errorf("%s:%s", op, "attemps count can't be 0")
	}

	var err error

	for range attemps {
		if err = fn(); err != nil {
			time.Sleep(delay)
			continue
		}
		return nil
	}

	return fmt.Errorf("%s:%w", op, err)
}
