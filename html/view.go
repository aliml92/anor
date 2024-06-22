package html

import (
	"embed"
	"fmt"
	"github.com/aliml92/anor/html/funcs"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

//go:embed templates/*
var templatesFS embed.FS

type View struct {
	templateCache map[string]*template.Template
}

func NewView() *View {
	templatesCache := make(map[string]*template.Template)
	err := parseFullPageTemplates(templatesCache)
	if err != nil {
		panic(err)
	}

	err = parseContentTemplates(templatesCache)
	if err != nil {
		panic(err)
	}

	err = parseComponentTemplates(templatesCache)
	if err != nil {
		panic(err)
	}

	return &View{
		templateCache: templatesCache,
	}

}

func (v *View) Render(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := v.templateCache[name]
	if !ok {
		http.Error(w, fmt.Sprintf("template %s not found", name), http.StatusInternalServerError)
		return
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to render template %s: %v", name, err), http.StatusInternalServerError)
	}
}

func (v *View) RenderComponent(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := v.templateCache[name]
	if !ok {
		http.Error(w, fmt.Sprintf("template %s not found", name), http.StatusInternalServerError)
		return
	}

	s := strings.Split(name, "/")
	name = strings.TrimSuffix(s[len(s)-1], ".gohtml")
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("template %s error: %v", name, err), http.StatusInternalServerError)
		return
	}
}

type Component map[string]any

func (v *View) RenderComponents(w http.ResponseWriter, components []Component) {
	for _, component := range components {
		for name, data := range component {
			tmpl, ok := v.templateCache[name]
			if !ok {
				http.Error(w, fmt.Sprintf("template %s not found", name), http.StatusInternalServerError)
				return
			}

			s := strings.Split(name, "/")
			name = strings.TrimSuffix(s[len(s)-1], ".gohtml")
			err := tmpl.ExecuteTemplate(w, name, data)
			if err != nil {
				http.Error(w, fmt.Sprintf("template %s error: %v", name, err), http.StatusInternalServerError)
				return
			}
		}

	}
}

func parseFullPageTemplates(templatesCache map[string]*template.Template) error {
	t, err := template.New("base").Funcs(funcs.FuncMap).ParseFS(templatesFS,
		"templates/base.gohtml",
		"templates/partials/header/*.gohtml",
		"templates/partials/*.gohtml",
	)
	if err != nil {
		return err
	}

	t, err = t.ParseFS(templatesFS, "templates/shared/*.gohtml")
	if err != nil {
		return err
	}

	dirEntries, err := templatesFS.ReadDir("templates/pages")
	if err != nil {
		return err
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			ti, err := t.Clone()
			if err != nil {
				return err
			}

			patterns := []string{
				"templates/pages/" + entry.Name() + "/content.gohtml",
				"templates/pages/" + entry.Name() + "/base.gohtml",
			}

			compDir, err := templatesFS.ReadDir("templates/pages/" + entry.Name() + "/components")
			if err == nil && len(compDir) > 0 {
				patterns = append(patterns, "templates/pages/"+entry.Name()+"/components/*.gohtml")
			}

			ti, err = ti.ParseFS(templatesFS, patterns...)
			if err != nil {
				return err
			}

			templatesCache["pages/"+entry.Name()+"/base.gohtml"] = ti
		}
	}

	return nil
}

func parseContentTemplates(templateCache map[string]*template.Template) error {
	dirEntries, err := templatesFS.ReadDir("templates/pages")
	if err != nil {
		return err
	}

	t, err := template.New("content").Funcs(funcs.FuncMap).ParseFS(templatesFS, "templates/shared/*.gohtml")
	if err != nil {
		return err
	}

	for _, entry := range dirEntries {
		ti, err := t.Clone()
		if err != nil {
			return err
		}

		compDir, err := templatesFS.ReadDir("templates/pages/" + entry.Name() + "/components")
		if err == nil && len(compDir) > 0 {
			// let's say `entry` is `cart` which is a folder name in templates/pages folder
			ti, err = ti.ParseFS(templatesFS, "templates/pages/"+entry.Name()+"/components/*.gohtml")
			if err != nil {
				return err
			}

		}

		ti, err = ti.ParseFS(templatesFS, "templates/pages/"+entry.Name()+"/content.gohtml")
		if err != nil {
			return err
		}

		templateCache["pages/"+entry.Name()+"/content.gohtml"] = ti
	}

	return nil
}

func parseComponentTemplates(templateCache map[string]*template.Template) error {
	err := parseComponentsInDir(templateCache, "templates/shared")
	if err != nil {
		return err
	}

	err = parseComponentsInDir(templateCache, "templates/partials/header")
	if err != nil {
		return err
	}

	dirEntries, err := templatesFS.ReadDir("templates/pages")
	if err != nil {
		return err
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			err := parseComponentsInDir(templateCache, "templates/pages/"+entry.Name()+"/components")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func parseComponentsInDir(templateCache map[string]*template.Template, dir string) error {
	dirEntries, err := templatesFS.ReadDir(dir)
	if err != nil {
		return nil
	}

	for _, entry := range dirEntries {
		ext := filepath.Ext(entry.Name())
		if ext == ".gohtml" {
			d := strings.TrimPrefix(dir, "templates/")
			templateCache[d+"/"+entry.Name()], err = template.New("").Funcs(funcs.FuncMap).ParseFS(templatesFS, filepath.Join(dir, entry.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
