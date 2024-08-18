package funcs

import (
	"bytes"
	"fmt"
	"github.com/aliml92/anor"
	"html/template"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

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

func FormatAddress(address anor.Address) template.HTML {
	var buffer bytes.Buffer
	lines := []struct {
		class string
		value string
	}{
		{"address-name", address.Name},
		{"address-line1", address.AddressLine1},
		{"address-line2", address.AddressLine2},
		{"", fmt.Sprintf(`<span class="address-city">%s</span>, <span class="address-state">%s</span> <span class="address-zip">%s</span>`,
			template.HTMLEscapeString(address.City),
			template.HTMLEscapeString(address.StateProvince),
			template.HTMLEscapeString(address.PostalCode))},
		{"address-country", address.Country},
	}

	for _, line := range lines {
		if strings.TrimSpace(line.value) != "" {
			if line.class != "" {
				buffer.WriteString(fmt.Sprintf(`<span class="%s">`, line.class))
				buffer.WriteString(template.HTMLEscapeString(line.value))
				buffer.WriteString("</span>")
			} else {
				buffer.WriteString(line.value) // For the city, state, zip line which is pre-formatted
			}
			buffer.WriteString("<br>")
		}
	}

	// Remove the last "<br>" if it exists
	html := buffer.String()
	if strings.HasSuffix(html, "<br>") {
		html = html[:len(html)-4]
	}

	return template.HTML(html)
}

func StepperClass(i, curr int) string {
	if i < curr {
		return "completed"
	} else if i == curr {
		return "active"
	}

	return ""
}

func AddressesEmpty(a, b anor.Address, c, d []anor.Address) bool {
	if a.IsZero() && b.IsZero() && len(c) == 0 && len(d) == 0 {
		return true
	}

	return false
}

// FormatDateTime formats a time.Time value for general date and time display
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

// FormatDate formats a time.Time value for date-only display
func FormatDate(t time.Time) string {
	return t.Format("02 January 2006")
}

// FormatTimeForJS formats a time.Time value for JavaScript parsing
func FormatTimeForJS(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05Z")
}

func ContainsPath(path, target string) bool {
	return strings.Contains(path, target)
}
