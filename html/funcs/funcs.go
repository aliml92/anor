package funcs

var FuncMap = map[string]interface{}{
	// basic arithmetic.
	"add":  Add,
	"add1": Add1,
	"sub":  Sub,
	"sub1": Sub1,
	"subd": Subd,
	"ned":  Ned,
	"div":  Div,
	"mul":  Mul,
	"muld": Muld,
	"mod":  Mod,

	// arithmetic comparisons
	"eqd0": Eqd0,

	"iterate": IterateInt32,

	"jsonify":        Jsonify,
	"currencySymbol": CurrencySymbol,

	// formatting
	"humanizeNum": HumanizeNum,
	"pluralize":   Pluralize,

	// application specific.
	"splitMap":                   SplitMap, // if size is odd then first half gets the remainder
	"modifyImgURL":               ModifyImgURL,
	"genPageNums":                GenPageNums,
	"formatHandle":               FormatHandle,
	"injectCategoryIntoSiblings": InjectCategoryIntoSiblings,
	"isBrandChecked":             IsBrandChecked,
	"formatProductQty":           FormatProductQty,
	"getRootCategoryAlias":       GetRootCategoryAlias,
	"formatAddress":              FormatAddress,
	"stepperClass":               StepperClass,
	"addressesEmpty":             AddressesEmpty,
	"formatDate":                 FormatDate,
	"formatDateTime":             FormatDateTime,
	"formatTimeForJS":            FormatTimeForJS,
	"containsPath":               ContainsPath,
}
