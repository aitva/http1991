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
	l = log.With(l, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	r := mux.NewRouter()
	{
		v := view.NewHTML("engine", "static/index")
		r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			v.Render(w, r, nil)
		})
	}
	{
		staticHandler := http.FileServer(http.Dir("./static/"))
		// staticHandler = http.StripPrefix("/static/", staticHandler)
		r.PathPrefix("/").Handler(staticHandler)
	}

	s := http.Server{
		Addr:    addr,
		Handler: r,
	}
	l.Log("msg", "listening on "+addr)
	s.ListenAndServe()
}
