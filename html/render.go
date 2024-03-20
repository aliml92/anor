package html

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/aliml92/anor"
	"github.com/aliml92/anor/html/funcs"
)

type Render struct {
	htmxTemplates *template.Template
	htmlTemplates map[string]*template.Template
}

func NewRender() *Render {
	hx := parseHTMXTemplates()

	h, err := parseHTMLTemplates()
	if err != nil {
		panic(err)
	}

	return &Render{
		htmxTemplates: hx,
		htmlTemplates: h,
	}
}

func (s *Render) HTML(w http.ResponseWriter, status int, page string, data interface{}) {
	ts, ok := s.htmlTemplates[page]
	if !ok {
		err := fmt.Errorf("the templates %s does not exist", page)
		slog.Error(err.Error())
		http.Error(w, anor.EINTERNALMSG, http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)

	if err := ts.ExecuteTemplate(buf, "base", data); err != nil {
		slog.Error(err.Error())
		http.Error(w, anor.EINTERNALMSG, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (s *Render) HTMX(w http.ResponseWriter, status int, page string, data interface{}) {
	buf := new(bytes.Buffer)

	if err := s.htmxTemplates.ExecuteTemplate(buf, page, data); err != nil {
		slog.Error(err.Error())
		http.Error(w, anor.EINTERNALMSG, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

func parseHTMLTemplates() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./html/templates/html/pages/*tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(funcs.FuncMap).ParseFiles("./html/templates/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./html/templates/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func parseHTMXTemplates() *template.Template {
	return template.Must(template.New("base").Funcs(funcs.FuncMap).ParseGlob("./html/templates/htmx/*.tmpl"))
}
