package funcs

import (
	"fmt"
	"regexp"

	"github.com/spf13/cast"
)

func ModifyImgURL(url string, resolution int) string {
	// Define a regular expression to match the last number in the URL
	re := regexp.MustCompile(`(\d+)(\.[a-zA-Z]+$)`)

	// Replace the last number with the new number
	newURL := re.ReplaceAllString(url, fmt.Sprintf("%d$2", resolution))

	return newURL
}

func SplitMap(m map[string]string) struct {
	First  map[string]string
	Second map[string]string
} {
	l := len(m)
	sLen := l / 2
	fLen := sLen + l%2

	i := 0
	first := make(map[string]string, fLen)
	second := make(map[string]string, sLen)
	for k, v := range m {
		if i < fLen {
			first[k] = v
		} else {
			second[k] = v
		}
		i++
	}

	return struct {
		First  map[string]string
		Second map[string]string
	}{first, second}
}

func HumanizeNum(num interface{}) string {
	str := cast.ToString(num)
	// Insert commas every three characters from the end
	for i := len(str) - 3; i > 0; i -= 3 {
		str = str[:i] + "," + str[i:]
	}

	return str
}

func GenPageNums(currentPage, pageCount, maxPages int) []int {
	s := maxPages
	if pageCount < maxPages {
		s = pageCount
	}

	c := currentPage
	if c > 2 {
		c -= 2
	} else if c > 1 {
		c--
	}

	res := make([]int, s)
	for i := 0; i < s; i++ {
		res[i] = c
		c++
	}

	return res
}

func FormatSlug(slug string, id interface{}) string {
	return fmt.Sprintf("%s-%v", slug, id)
}
