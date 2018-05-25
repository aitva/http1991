package view

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-kit/kit/log"
)

// Default value for the package.
const (
	LayoutDir   = "view/layout/"
	TemplateDir = "view/"
	TemplateExt = ".gohtml"
)

// HTMLConfig contains optional configuration for the HTML view.
type HTMLConfig func(*HTML)

// HTMLSetLogger sets a logger for the view.
func HTMLSetLogger(l log.Logger) HTMLConfig {
	return func(html *HTML) {
		html.l = l
	}
}

// HTML represents an HTML view.
type HTML struct {
	Template *template.Template
	Layout   string
	files    []string
	l        log.Logger
}

// NewHTML instanciates a new HTML view using a given file
// and layout.
func NewHTML(layout string, files []string, opts ...HTMLConfig) *HTML {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)
	t, err := template.New("").ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	h := &HTML{
		Template: t,
		Layout:   layout,
		files:    files,
		l:        log.NewNopLogger(),
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// Render renders an HTML view onto an HTTP response.
func (h *HTML) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	var buf bytes.Buffer

	err := h.Template.ExecuteTemplate(&buf, h.Layout, data)
	if err != nil {
		h.l.Log("err", err)
		http.Error(w, `Something went wrong. If the problem persists,
please email arod.louis@gmail.com`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	io.Copy(w, &buf)
}

// HTMLWatcher watches a set of views.
type HTMLWatcher struct {
	watcher *fsnotify.Watcher
	files   map[string][]*HTML
	l       log.Logger
}

// NewHTMLWatcher instanciates a new HTMLWatcher.
func NewHTMLWatcher(l log.Logger, htmls ...*HTML) *HTMLWatcher {
	hw := &HTMLWatcher{}

	hw.l = log.NewNopLogger()
	if l != nil {
		hw.l = l
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	hw.watcher = w

	files := make(map[string][]*HTML)
	for _, h := range htmls {
		for _, f := range h.files {
			htmls := files[f]
			htmls = append(htmls, h)
			files[f] = htmls
		}
	}
	hw.files = files

	for f := range hw.files {
		err := hw.watcher.Add(f)
		if err != nil {
			panic(err)
		}
	}
	return hw
}

// Watch watches and reload HTML views.
func (hw *HTMLWatcher) Watch() {
	// we use files and ticker to dedup event
	files := map[string]struct{}{}
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	w := hw.watcher
	for {
		select {
		case <-t.C:
			for file := range files {
				hw.reload(file)
			}
			if len(files) > 0 {
				files = map[string]struct{}{}
			}
		case e, ok := <-w.Events:
			if !ok {
				return
			}
			files[e.Name] = struct{}{}
		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			hw.l.Log("err", err)
		}
	}
}

func (hw *HTMLWatcher) reload(file string) {
	htmls := hw.files[file]
	for _, h := range htmls {
		t, err := template.New("").ParseFiles(h.files...)
		if err != nil {
			hw.l.Log("msg", "fail to reload template", "err", err)
			break
		}
		h.Template = t
	}
}

// Close closes the watcher.
func (hw *HTMLWatcher) Close() error {
	return hw.Close()
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
