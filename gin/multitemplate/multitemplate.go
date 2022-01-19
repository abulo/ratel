package multitemplate

import (
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"

	"github.com/abulo/ratel/v2/gin/render"
)

// Render type
type Render map[string]*template.Template

var (
	_ render.HTMLRender = Render{}
	_ Renderer          = Render{}
)

// New instance
func New() Render {
	return make(Render)
}

// Add new template
func (r Render) Add(name string, tmpl *template.Template) {
	if tmpl == nil {
		panic("template can not be nil")
	}
	if len(name) == 0 {
		panic("template name cannot be empty")
	}
	if _, ok := r[name]; ok {
		panic(fmt.Sprintf("template %s already exists", name))
	}
	r[name] = tmpl
}

// AddFromFiles supply add template from files
func (r Render) AddFromFiles(name string, files ...string) *template.Template {
	tmpl := template.Must(template.ParseFiles(files...))
	r.Add(name, tmpl)
	return tmpl
}

// AddFromGlob supply add template from global path
func (r Render) AddFromGlob(name, glob string) *template.Template {
	tmpl := template.Must(template.ParseGlob(glob))
	r.Add(name, tmpl)
	return tmpl
}

// AddFromString supply add template from strings
func (r Render) AddFromString(name, templateString string) *template.Template {
	tmpl := template.Must(template.New(name).Parse(templateString))
	r.Add(name, tmpl)
	return tmpl
}

// AddFromStringsFuncs supply add template from strings
func (r Render) AddFromStringsFuncs(name string, funcMap template.FuncMap, templateStrings ...string) *template.Template {
	tmpl := template.New(name).Funcs(funcMap)

	for _, ts := range templateStrings {
		tmpl = template.Must(tmpl.Parse(ts))
	}

	r.Add(name, tmpl)
	return tmpl
}

// AddFromFilesFuncs supply add template from file callback func
func (r Render) AddFromFilesFuncs(name string, funcMap template.FuncMap, files ...string) *template.Template {
	tname := filepath.Base(files[0])
	tmpl := template.Must(template.New(tname).Funcs(funcMap).ParseFiles(files...))
	r.Add(name, tmpl)
	return tmpl
}

func (r Render) AddFromFs(name string, fs fs.FS, patterns ...string) *template.Template {
	tmpl := template.Must(template.ParseFS(fs, patterns...))
	r.Add(name, tmpl)
	return tmpl
}

func (r Render) AddFromFsFuncs(name string, funcMap template.FuncMap, fs fs.FS, patterns ...string) *template.Template {
	tname := filepath.Base(patterns[0])
	tmpl := template.Must(template.New(tname).Funcs(funcMap).ParseFS(fs, patterns...))
	r.Add(name, tmpl)
	return tmpl
}

// Instance supply render string
func (r Render) Instance(name string, data interface{}) render.Render {
	return render.HTML{
		Template: r[name],
		Data:     data,
	}
}
