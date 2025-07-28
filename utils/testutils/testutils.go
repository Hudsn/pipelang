package testutils

import (
	"fmt"
)

// returns whether the two items are equal, and a string to use in the test error message if the resulting bool is undesired.
func Equal[T comparable](check T, got T) (bool, string) {
	if check == got {
		return true, fmt.Sprintf("wrong value. values should not be equal:\n\tdon't want=%+v\n\tgot=%+v", check, got)
	} else {
		return false, fmt.Sprintf("wrong value. values should be equal:\n\twant=%+v\n\tgot=%+v", check, got)
	}
}
