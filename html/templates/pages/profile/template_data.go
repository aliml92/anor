package profile

type Content struct {
}

func (Content) GetTemplateFilename() string {
	return "content.gohtml"
}
