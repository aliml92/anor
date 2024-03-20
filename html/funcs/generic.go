package funcs

import (
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

// Basic arithmetic functions

func Add(i ...interface{}) int64 {
	var a int64 = 0
	for _, b := range i {
		a += toInt64(b)
	}
	return a
}

func Sub(a, b interface{}) int64 {
	return toInt64(a) - toInt64(b)
}

func Subd(a, b decimal.Decimal) decimal.Decimal {
	return a.Sub(b)
}

func Eqd0(a decimal.Decimal) bool {
	return a.Equal(decimal.Zero)
}

func Div(a, b interface{}) int64 {
	return toInt64(a) / toInt64(b)
}

func Mod(a, b interface{}) int64 {
	return toInt64(a) % toInt64(b)
}

func Add1(i interface{}) int64 {
	return toInt64(i) + 1
}

func Sub1(i interface{}) int64 {
	return toInt64(i)
}

func toInt64(i interface{}) int64 {
	return cast.ToInt64(i)
}
