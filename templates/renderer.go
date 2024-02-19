package templates

import (
	"html/template"
	"io"
	"time"

	"github.com/labstack/echo"
	"github.com/ladydascalie/quakes/config/locales"
)

// Tpls is for global access to templates
var Tpls *Template

// TimeLayout defines the standard display format for time inside templates
var TimeLayout = "02 Jan 06 15:04:05"

// New is used here to glob the template files and add the function map to the templates
func New() *Template {
	funcMap := template.FuncMap{
		"T": func(msg string) string {
			return locales.T(msg)
		},
		"formatTime": func(t time.Time) string {
			return t.Format(TimeLayout)
		},
	}
	Tpls = &Template{
		Templates: template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html")),
	}
	return Tpls
}

// Template wraps template.Template
type Template struct {
	Templates *template.Template
}

// Render wraps templates.ExecuteTemplate.
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}
