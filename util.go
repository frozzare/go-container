package container

import (
	"regexp"
	"strings"
)

func getFuncNameParameters(name string) []string {
	r := regexp.MustCompile(`\((.*)\)`)
	o := r.FindStringSubmatch(name)

	if len(o) < 1 {
		return []string{}
	}

	o[1] = strings.Replace(o[1], " ", "", -1)

	return strings.Split(o[1], ",")
}

func getFieldTypeName(name string) string {
	r := regexp.MustCompile(`\<(.*)\s`)
	o := r.FindStringSubmatch(name)

	if len(o) < 1 {
		return ""
	}

	return strings.Replace(o[1], " ", "", -1)
}
