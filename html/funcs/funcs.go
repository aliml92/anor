package funcs

var FuncMap = map[string]interface{}{
	// basic arithmetic.
	"add":  Add,
	"add1": Add1,
	"sub":  Sub,
	"sub1": Sub1,
	"subd": Subd,
	"div":  Div,
	"mod":  Mod,

	// arithmetic comparisons
	"eqd0": Eqd0,

	// formatting
	"humanizeNum": HumanizeNum,

	// application specific.
	"splitMap":     SplitMap, // if size is odd then first half gets the remainder
	"modifyImgURL": ModifyImgURL,
	"genPageNums":  GenPageNums,
	"formatSlug":   FormatSlug,
}
