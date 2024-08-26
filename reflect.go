package core

import "strings"

// ParseTags parses a string of tags into a map.
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
			res[tag] = "true"
		} else {
			res[frames[0]] = frames[1]
		}
	}

	return res
}
