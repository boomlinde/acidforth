package machine

import "regexp"

func TokenizeString(source string) []string {
	matcher := regexp.MustCompile(`\S+`)
	res := matcher.FindAllString(source, -1)
	if res == nil {
		return make([]string, 0)
	}
	return res
}

func TokenizeBytes(source []byte) []string {
	matcher := regexp.MustCompile(`\S+`)
	res := matcher.FindAll(source, -1)
	if res == nil {
		return make([]string, 0)
	}
	strings := make([]string, len(res))
	for i, v := range res {
		strings[i] = string(v)
	}
	return strings
}
