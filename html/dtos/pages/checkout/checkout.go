package checkout

import (
	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html/dtos/partials"
)

type Base struct {
	partials.Header
	Content
}

type Content struct {
	Cart anor.Cart
}
