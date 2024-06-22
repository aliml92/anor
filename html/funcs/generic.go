package funcs

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"html/template"
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

func Ned(a, b decimal.Decimal) bool {
	return !a.Equal(b)
}

func Eqd0(a decimal.Decimal) bool {
	return a.Equal(decimal.Zero)
}

func Div(a, b interface{}) int64 {
	return toInt64(a) / toInt64(b)
}

func Mul(a, b interface{}) int64 {
	return toInt64(a) * toInt64(b)
}

func Muld(a decimal.Decimal, b int32) decimal.Decimal {
	x := decimal.NewFromInt32(b)
	return a.Mul(x)
}

func Mod(a, b interface{}) int64 {
	return toInt64(a) % toInt64(b)
}

func Add1(i interface{}) int64 {
	return toInt64(i) + 1
}

func Sub1(i interface{}) int64 {
	return toInt64(i) - 1
}

func toInt64(i interface{}) int64 {
	return cast.ToInt64(i)
}

func IterateInt32(max int32) []int32 {
	if max < 1 {
		return nil
	}

	if max > 10 {
		max = 10
	}

	res := make([]int32, max)
	for i := range res {
		res[i] = int32(i) + 1
	}

	return res
}

func Jsonify(data interface{}) template.JS {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return template.JS(jsonBytes)
}

func CurrencySymbol(code string) string {
	switch code {
	case "USD":
		return "$"
	case "EUR":
		return "€"
	case "GBP":
		return "£"
	case "JPY":
		return "¥"
	case "AUD":
		return "A$"
	case "CAD":
		return "C$"
	case "CHF":
		return "Fr"
	case "CNY":
		return "¥"
	case "HKD":
		return "HK$"
	case "NZD":
		return "NZ$"
	case "SEK":
		return "kr"
	case "KRW":
		return "₩"
	case "SGD":
		return "S$"
	case "NOK":
		return "kr"
	case "MXN":
		return "M$"
	case "INR":
		return "₹"
	case "RUB":
		return "₽"
	case "TRY":
		return "₺"
	case "BRL":
		return "R$"
	default:
		return code
	}
}
