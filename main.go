package main

import (
	"net/http"
	"os"

	"github.com/codex-veritas/http1991/view"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
)

const addr = ":1991"

func main() {
	w := log.NewSyncWriter(os.Stdout)
	l := log.NewLogfmtLogger(w)
	l = log.With(l, "ts", log.DefaultTimestamp, "caller", log.DefaultCaller)

	r := mux.NewRouter()
	r.StrictSlash(true)
	views := []*view.HTML{}
	{
		v := view.NewHTML(
			"default",
			[]string{"static/index"},
			view.HTMLSetLogger(l),
		)
		r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			v.Render(w, r, view.Engine{Title: "Evil 1991"})
		})
		views = append(views, v)
	}
	{
		v := view.NewHTML(
			"engine",
			[]string{"game/intern-office"},
			view.HTMLSetLogger(l),
		)
		r.HandleFunc("/game/intern-office/", func(w http.ResponseWriter, r *http.Request) {
			v.Render(w, r, view.Engine{Title: "Intern Office"})
		})
		views = append(views, v)
	}
	{
		staticHandler := http.FileServer(http.Dir("./static/"))
		r.PathPrefix("/").Handler(staticHandler)
	}

	hw := view.NewHTMLWatcher(l, views...)
	defer hw.Close()
	go hw.Watch()

	s := http.Server{
		Addr:    addr,
		Handler: r,
	}
	l.Log("msg", "listening on "+addr)
	s.ListenAndServe()
}
