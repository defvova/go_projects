package service

import (
	"html/template"
	"net/http"
	"strings"
	"time"
)

type Renderer struct {
	Base  *template.Template
	Pages map[string]*template.Template
}

var funcMap = template.FuncMap{
	"upper": strings.ToUpper,
	"truncate": func(s string, n int) string {
		if len(s) <= n {
			return s
		}
		return s[:n] + "..."
	},
	"fmtDate": func(t time.Time) string {
		return t.Format("02 Jan 2006")
	},
}

func NewRenderer(pageNames []string) *Renderer {
	base := template.Must(template.New("base").Funcs(funcMap).ParseFiles(
		"web/layouts/base.html",
		"web/shared/nav.html",
		"web/shared/alert.html",
	))

	pages := map[string]*template.Template{}
	for _, name := range pageNames {
		t := template.Must(template.Must(base.Clone()).ParseFiles("web/pages/" + name + ".html"))
		pages[name] = t
	}

	return &Renderer{Base: base, Pages: pages}
}

func (r *Renderer) Render(w http.ResponseWriter, pageName string, data any) {
	t, ok := r.Pages[pageName]
	if !ok {
		http.Error(w, "template not found", http.StatusNotFound)
		return
	}
	if err := t.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
