package util

import "errors"

func IsErrors(err error, errs []error) bool {
	for _, e := range errs {
		if errors.Is(err, e) {
			return true
		}
	}
	return false
}
