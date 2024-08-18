package signup

type Content struct{}

func (Content) GetTemplateFilename() string {
	return "content.gohtml"
}
