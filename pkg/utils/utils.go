package utils

import (
	"math/rand"
	"regexp"
	"strings"

	gonanoid "github.com/matoous/go-nanoid"
	"golang.org/x/crypto/bcrypt"
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

func CreateSlug(val string) string {
	// Replace non-alphanumeric characters (except comma: ,) with a hyphen
	slug := strings.ReplaceAll(val, "&", "and")

	reg := regexp.MustCompile("[^A-Za-z0-9,]+")
	slug = reg.ReplaceAllString(slug, "-")

	// Replace all consecutive hyphens with single hyphen
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	// check if there is a comma between 60 and 100th positions
	var commaIdx int
	if len(slug) > 60 {
		for i := 60; i < len(slug) && i < 100; i++ {
			if slug[i] == ',' {
				commaIdx = i
				break
			}
		}
	}

	if commaIdx != 0 {
		// slice the slug until the pos of comma
		slug = slug[:commaIdx]
	} else {
		// Truncate if length exceeds 80 characters
		maxLength := 80
		if len(slug) > maxLength {
			// strip out the slug until the pos of maxLength
			slug = slug[:maxLength]

			// Remove the last partial word if exists
			slug = slug[:strings.LastIndexByte(slug, '-')]
		}
	}

	// remove remaining colons
	slug = strings.ReplaceAll(slug, ",", "")
	// Convert to lowercase
	slug = strings.ToLower(slug)

	return slug
}
