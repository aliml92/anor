package funcs

import (
	"fmt"
	"github.com/aliml92/anor"
	"regexp"
	"sort"
	"strconv"

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
	if c > 3 {
		c -= 3
	} else if c > 2 {
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

func FormatHandle(handle string, id interface{}) string {
	return fmt.Sprintf("%s-%v", handle, id)
}

func InjectCategoryIntoSiblings(c anor.Category, siblings []anor.Category) []anor.Category {
	var categories []anor.Category
	if len(siblings) < 15 {
		categories = append(siblings, c)
		sort.Slice(categories, func(i, j int) bool {
			return categories[i].ID < categories[j].ID
		})
	} else {
		categories = append([]anor.Category{c}, siblings...)
	}

	return categories
}

func IsBrandChecked(brand string, brands []string) bool {
	if len(brands) == 0 {
		return false
	}
	for _, b := range brands {
		if b == brand {
			return true
		}
	}
	return false
}

func FormatProductQty(count int) string {
	if count < 1 {
		return "None left"
	}

	if count == 1 {
		return "Only one left"
	}

	if count > 10 {
		return "More than 10 available"
	}

	return strconv.Itoa(count) + " left"
}

func GetRootCategoryAlias(category string) string {
	switch category {
	case "Clothing, Shoes & Accessories":
		return "Fashion"
	case "Jewelry & Watches":
		return "Jewelry"
	default:
		return category
	}
}
