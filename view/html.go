package view

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/go-kit/kit/log"
)

const (
	TemplateDir = "views/layout"
	TemplateExt = ".gohtml"
)

type HTML struct {
	Template *template.Template
	Layout   string
	l        *log.Logger
}

func NewHTML(layout string, files ...string) *HTML {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)
	t, err := template.New("").ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
		l:        log.NewNopLogger(),
	}

}

func (h *HTML) SetLogger(l *log.Logger) {
	h.l = l
}

func (h *HTML) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	var buf bytes.Buffer

	err := h.Template.ExecuteTemplate(&buf, h.Layout, nil)
	if err != nil {
		l.Log("err", err)
		http.Error(w, `Something went wrong. If the problem persists,
please email arod.louis@gmail.com`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	io.Copy(w, &buf)
}

// layoutFiles returns a slice of strings representing
// the layout files used in our application.
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

// addTemplatePath takes in a slice of strings
// representing file paths for templates, and it prepends
// the TemplateDir directory to each string in the slice
//
// Eg the input {"home"} would result in the output
// {"views/home"} if TemplateDir == "views/"
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

// addTemplateExt takes in a slice of strings
// representing file paths for templates and it appends
// the TemplateExt extension to each string in the slice
//
// Eg the input {"home"} would result in the output
// {"home.gohtml"} if TemplateExt == ".gohtml"
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
