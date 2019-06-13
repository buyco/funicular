package utils

import "reflect"

func InArray(needle interface{}, haystack interface{}) (bool, int) {
	var exists = false
	var index = -1

	switch reflect.TypeOf(haystack).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(haystack)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(needle, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return exists, index
			}
		}
	}
	return exists, index
}
