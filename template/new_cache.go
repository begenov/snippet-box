package templates

import (
	"html/template"
	"path/filepath"
	"snippet/internal/service/forms"
	"snippet/model"
)

type TemplateData struct {
	CurrentYear     int
	Form            *forms.Form
	Snippet         *model.Snippet
	Snippets        []*model.Snippet
	IsAuthenticated bool
	Flash           string
}

func NewTemplatesData(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}
	for _, page := range pages {

		name := filepath.Base(page)
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partioal.html"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
