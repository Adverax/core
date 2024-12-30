package core

import (
	"reflect"
	"strings"
)

func IsZeroValue(x interface{}) bool {
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func ParseTags(tags string) map[string]string {
	list := strings.Split(tags, ",")
	res := make(map[string]string)
	for i, tag := range list {
		if tag == "" {
			continue
		}

		var frames []string
		frames = strings.Split(tag, "=")
		if i == 0 {
			if len(frames) == 1 {
				res["name"] = frames[0]
				continue
			}
		}

		if len(frames) == 1 {
			res[tag] = ""
		} else {
			res[frames[0]] = frames[1]
		}
	}

	return res
}
