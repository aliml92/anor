package reset_password

type Content struct{}

func (Content) GetTemplateFilename() string {
	return "content.gohtml"
}
