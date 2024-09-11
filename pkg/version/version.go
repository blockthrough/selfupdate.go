package version

import (
	"regexp"
	"strings"
)

var onlyNumberRegExp = regexp.MustCompile(`[^0-9]`)

// if a > b return true
func Compare(a, b string) bool {
	listA := strings.Split(a, ".")
	listB := strings.Split(b, ".")

	// adding 0 to the shorter list
	if len(listA) > len(listB) {
		for i := 0; i < len(listA)-len(listB); i++ {
			listB = append(listB, "0")
		}
	} else {
		for i := 0; i < len(listB)-len(listA); i++ {
			listA = append(listA, "0")
		}
	}

	// make sure each element is a number
	for i := 0; i < len(listA); i++ {
		listA[i] = onlyNumberRegExp.ReplaceAllString(listA[i], "")
	}

	for i := 0; i < len(listB); i++ {
		listB[i] = onlyNumberRegExp.ReplaceAllString(listB[i], "")
	}

	// adding leading 0 to each element
	for i := 0; i < len(listA); i++ {
		if len(listA[i]) < len(listB[i]) {
			listA[i] = strings.Repeat("0", len(listB[i])-len(listA[i])) + listA[i]
		} else if len(listA[i]) > len(listB[i]) {
			listB[i] = strings.Repeat("0", len(listA[i])-len(listB[i])) + listB[i]
		}
	}

	for i := 0; i < len(listA); i++ {
		if listA[i] == listB[i] {
			continue
		}

		return listA[i] > listB[i]
	}

	return len(listA) > len(listB)
}
