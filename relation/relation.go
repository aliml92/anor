package relation

type Relation string

const (
	OrderItems        Relation = "OrderItems"
	StripeCardPayment Relation = "StripeCardPayment"
	ShippingAddress   Relation = "ShippingAddress"
	BillingAddress    Relation = "BillingAddress"
	User              Relation = "User"
	Category          Relation = "Category"
	Product           Relation = "Product"
	Store             Relation = "Store"
	Pricing           Relation = "Pricing"
	ProductVariant    Relation = "ProductVariant"
	ProductAttribute  Relation = "ProductAttribute"
	Rating            Relation = "Rating"
	Review            Relation = "Review"
)

type Set map[Relation]struct{}

func New(relations ...Relation) Set {
	set := make(Set)
	for _, relation := range relations {
		set[relation] = struct{}{}
	}
	return set
}

func (s Set) Has(relation Relation) bool {
	_, ok := s[relation]
	return ok
}
