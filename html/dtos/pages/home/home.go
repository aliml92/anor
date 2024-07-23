package home

import (
	"github.com/aliml92/anor/html/dtos/pages/home/components"
	"github.com/aliml92/anor/html/dtos/partials"
)

type Base struct {
	partials.Header
	Content
}

type Content struct {
	Featured    components.Featured
	NewArrivals components.Collection
	Popular     components.Collection
}
