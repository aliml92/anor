package utils

import (
	gonanoid "github.com/matoous/go-nanoid"
	sw "github.com/toadharvard/stopwords-iso"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

var discounts = []float32{
	0, 0, 0.02, 0, 0, 0, 0.03, 0, 0.04, 0, 0, 0.05, 0, 0, 0.07, 0, 0, 0, 0.08, 0, 0, 0, 0.09, 0,
	0.10, 0, 0.11, 0, 0.12, 0, 0.13, 0, 0.14, 0, 0.15, 0, 0.20, 0, 0, 0.25, 0, 0, 0, 0, 0, 0, 0,
	0.30, 0, 0, 0.35, 0, 0.40, 0, 0.45, 0, 0.50, 0.55, 0, 0.60, 0, 0.65, 0, 0.70, 0, 0, 0, 0, 0,
	0.15, 0, 0.12, 0, 0, 0, 0.05, 0.07, 0.20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

func GeneraterRandomCode(n int) string {
	const digits = "0123456789"
	otp := make([]byte, n)

	for i := 0; i < n; i++ {
		randomDigit := digits[rand.Intn(len(digits))]
		otp[i] = randomDigit
	}

	return string(otp)
}

func GenRandomNum(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func GenerateRandomDiscount() float32 {
	idx := rand.Intn(len(discounts))
	return discounts[idx]
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func MustGenerateID(l int) string {
	return gonanoid.MustGenerate("0123456789", l)
}

func CreateHandle(val string) string {
	// Replace non-alphanumeric characters (except comma: ,) with a hyphen
	handle := strings.ReplaceAll(val, "&", "and")

	reg := regexp.MustCompile("[^A-Za-z0-9,]+")
	handle = reg.ReplaceAllString(handle, "-")

	// Replace all consecutive hyphens with single hyphen
	reg = regexp.MustCompile("-+")
	handle = reg.ReplaceAllString(handle, "-")

	// Remove leading and trailing hyphens
	handle = strings.Trim(handle, "-")

	// check if there is a comma between 60 and 100th positions
	var commaIdx int
	if len(handle) > 60 {
		for i := 60; i < len(handle) && i < 100; i++ {
			if handle[i] == ',' {
				commaIdx = i
				break
			}
		}
	}

	if commaIdx != 0 {
		// slice the handle until the pos of comma
		handle = handle[:commaIdx]
	} else {
		// Truncate if length exceeds 80 characters
		maxLength := 80
		if len(handle) > maxLength {
			// strip out the handle until the pos of maxLength
			handle = handle[:maxLength]

			// Remove the last partial word if exists
			handle = handle[:strings.LastIndexByte(handle, '-')]
		}
	}

	// remove remaining colons
	handle = strings.ReplaceAll(handle, ",", "")
	// Convert to lowercase
	handle = strings.ToLower(handle)

	return handle
}

type SKUGenerator struct {
	base string
	p    string
	c    int

	varSkuSuffices map[string]struct{}
}

func NewSKUGenerator(category, product string) *SKUGenerator {
	swr, _ := sw.NewStopwordsMapping()
	bs := generateBaseSKU(category, product, swr)
	bs = strings.ToUpper(bs)
	return &SKUGenerator{
		base:           bs,
		p:              "S",
		c:              1,
		varSkuSuffices: make(map[string]struct{}),
	}
}

func (s *SKUGenerator) GetBaseSKU() string {
	return s.base
}

func (s *SKUGenerator) GenerateSKU(attributes []string) string {
	if len(attributes) == 0 {
		return s.base
	}

	var b strings.Builder
	b.WriteString(s.base) // e.g. "CLO34-NEKO3264"

	skuSuffix := s.genSKUSuffix(attributes)
	_, ok := s.varSkuSuffices[skuSuffix]
	if ok {
		s.c = s.c + 1
		skuSuffix = s.genSKUSuffix(attributes)
	}
	s.varSkuSuffices[skuSuffix] = struct{}{}

	b.WriteString("-")
	b.WriteString(skuSuffix)

	return b.String()
}

func (s *SKUGenerator) genSKUSuffix(attributes []string) string {
	var b strings.Builder
	b.WriteString(s.p)
	cStr := strconv.Itoa(s.c)
	b.WriteString(cStr)
	b.WriteString("-")

	l := len(attributes)
	for index, attr := range attributes {
		attr = strings.TrimSpace(attr)
		attr = strings.ToUpper(attr)
		re := regexp.MustCompile("[^a-zA-Z0-9]")
		attr = re.ReplaceAllString(attr, "")
		if len(attr) >= 2 {
			b.WriteString(string(attr[:2]))
		} else {
			b.WriteString(attr)
		}
		if index < l-1 {
			b.WriteString("/")
		}
	}

	return b.String()
}

func generateBaseSKU(category string, productName string, swr sw.StopwordsMapping) string {
	var builder strings.Builder

	ci := getCategoryInitials(category)
	builder.WriteString(ci)

	cs := generateCategorySuffix()
	builder.WriteString(cs)

	builder.WriteString("-")

	productName = swr.ClearString(productName)
	productWords := strings.Fields(productName)
	if len(productWords) > 1 {
		var wCount = 0
		for _, word := range productWords {
			if len(word) <= 2 {
				continue
			}
			if wCount >= 2 {
				break
			}
			builder.WriteString(word[:2])
			wCount++
		}
	} else {
		if len(productName) >= 4 {
			builder.WriteString(productName[:4])
		} else {
			builder.WriteString(productName)
		}
	}

	ps := generateProductSuffix()
	builder.WriteString(ps)

	return builder.String()
}

func getCategoryInitials(category string) string {
	// Remove non-English letters using regex
	re := regexp.MustCompile("[^a-zA-Z ]")
	category = re.ReplaceAllString(category, "")

	// Split category into individual words
	words := strings.Fields(category)

	// Initialize result variable
	var result string

	// Process based on the number of words
	switch len(words) {
	case 1:
		// If one word, get the first three letters
		result = words[0][:3]
	case 2:
		// If two words, get first two letters from first word and first letter from second word
		result = words[0][:2] + string(words[1][0])
	case 3:
		// If three words, get first letter from each word
		for _, word := range words {
			result += string(word[0])
		}
	default:
		// If more than three words, only consider the first three
		for i := 0; i < 3; i++ {
			result += string(words[i][0])
		}
	}

	// Convert result to uppercase
	result = strings.ToUpper(result)

	return result
}

func generateCategorySuffix() string {
	randomNum := GenRandomNum(10, 99)
	return strconv.Itoa(randomNum)
}

func generateProductSuffix() string {
	randomNum := GenRandomNum(1000, 9999)
	return strconv.Itoa(randomNum)
}
