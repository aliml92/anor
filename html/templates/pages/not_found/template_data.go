package not_found

type Content struct {
	Message string
}

func (Content) GetTemplateFilename() string {
	return "content.gohtml"
}
